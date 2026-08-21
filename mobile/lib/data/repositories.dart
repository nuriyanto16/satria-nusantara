import 'dart:convert';
import 'package:dio/dio.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../core/network.dart';
import '../core/constants.dart';
import 'models.dart';

class AuthRepository {

  Future<void> registerUser({
    required String email,
    required String password,
    required String name,
    required String phone,
    required String unit,
    required String tingkat,
    String? birthDate,
    String? gender,
    String? googleId,
  }) async {
    final cleanEmail = email.trim().toLowerCase();

    try {
      await api.dio.post(
        '/auth/signup-anggota',
        data: {
          'email': cleanEmail,
          'password': password,
          'nama_lengkap': name,
          'no_hp': phone,
          'unit_id': unit,
          'tingkatan': tingkat,
          if (birthDate != null && birthDate.isNotEmpty) 'tanggal_lahir': birthDate,
          if (gender != null && gender.isNotEmpty) 'jenis_kelamin': gender,
          if (googleId != null && googleId.isNotEmpty) 'google_id': googleId,
        },
      );
    } catch (e) {
      if (e is DioException) {
        final errMsg = e.response?.data['message'] ?? e.message ?? e.toString();
        throw Exception(errMsg);
      }
      rethrow;
    }
  }

  void approveUser(String email) {
    // Left empty or can call backend to approve
  }

  Future<void> logout() async {
    try {
      final googleSignIn = GoogleSignIn();
      if (await googleSignIn.isSignedIn()) {
        await googleSignIn.disconnect();
        await googleSignIn.signOut();
      }
    } catch (_) {}
    api.setToken(null);
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('auth_token');
    await prefs.remove('auth_user');
  }

  Future<Map<String, dynamic>?> checkSession() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final token = prefs.getString('auth_token');
      final userStr = prefs.getString('auth_user');

      if (token != null && token.isNotEmpty && userStr != null && userStr.isNotEmpty) {
        api.setToken(token);
        final user = User.fromJson(jsonDecode(userStr));
        return {'token': token, 'user': user};
      }
    } catch (e) {
      // Ignore
    }
    return null;
  }

  Future<Map<String, dynamic>> login(String email, String password) async {
    final cleanEmail = email.trim().toLowerCase();

    try {
      // Call to real backend API
      final response = await api.dio.post(
        ApiConstants.login,
        data: {'email': email, 'password': password},
      );
      final data = response.data['data'];
      final token = data['token'];
      final user = User.fromJson(data['user']);
      
      // Save token in API client
      api.setToken(token);
      
      // Save to SharedPreferences
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('auth_token', token);
      await prefs.setString('auth_user', jsonEncode(data['user']));
      
      return {'token': token, 'user': user};
    } catch (e) {
      if (e is DioException) {
        if (e.response?.statusCode == 403) {
          throw Exception("PENDING_VERIFICATION");
        }
        final errMsg = e.response?.data['message'] ?? "Email atau password salah";
        throw Exception(errMsg);
      }
      rethrow;
    }
  }

  Future<Map<String, dynamic>> loginGoogle(String email, String name, {String? googleId, String? noHp, String? fotoUrl, String? unitId}) async {
    final cleanEmail = email.trim().toLowerCase();
    try {
      final response = await api.dio.post(
        ApiConstants.googleLogin,
        data: {
          'email': cleanEmail,
          'nama_lengkap': name,
          if (googleId != null && googleId.isNotEmpty) 'google_id': googleId,
          if (noHp != null && noHp.isNotEmpty) 'no_hp': noHp,
          if (fotoUrl != null && fotoUrl.isNotEmpty) 'foto_url': fotoUrl,
          if (unitId != null && unitId.isNotEmpty) 'unit_id': unitId,
        },
      );
      final data = response.data['data'];
      final token = data['token'];
      final user = User.fromJson(data['user']);
      
      api.setToken(token);
      
      // Save to SharedPreferences
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('auth_token', token);
      await prefs.setString('auth_user', jsonEncode(data['user']));

      return {'token': token, 'user': user};
    } catch (e) {
      if (e is DioException) {
        if (e.response?.statusCode == 403) {
          throw Exception("PENDING_VERIFICATION");
        }
        if (e.response?.statusCode == 404) {
          rethrow;
        }
      }

      // Offline / Fallback handling for prototype testing removed
      // Force profile completion if not found
      throw Exception("User not found");
    }
  }

}

class SesiRepository {
  Future<List<Sesi>> getSesiList({String? unitId, String? tanggal}) async {
    try {
      final response = await api.dio.get(
        ApiConstants.getSesi,
        queryParameters: {
          if (unitId != null && unitId.isNotEmpty) 'unit_id': unitId,
          if (tanggal != null && tanggal.isNotEmpty) 'tanggal': tanggal,
        },
      );
      final List<dynamic> list = response.data['data'] ?? [];
      return list.map((item) => Sesi.fromJson(item)).toList();
    } catch (e) {
      if (e is DioException) {
        throw Exception(e.response?.data['message'] ?? e.message ?? e.toString());
      }
      rethrow;
    }
  }

  Future<void> scanQR(String code, String userId) async {
    await api.dio.post(
      ApiConstants.scanQr,
      data: {
        'qr_code': code,
        'user_id': userId,
      },
    );
  }
}

class FinanceRepository {
  Future<List<Iuran>> getIuranHistory(String userId) async {
    try {
      final response = await api.dio.get(
        '/finance/iuran',
        queryParameters: {'userId': userId},
      );
      final List<dynamic> list = response.data['data'] ?? [];
      return list.map((item) => Iuran.fromJson(item)).toList();
    } catch (e) {
      if (e is DioException) {
        throw Exception(e.response?.data['message'] ?? e.message ?? e.toString());
      }
      rethrow;
    }
  }

  Future<void> addMockIuran(String userId, {int? bulan, int? tahun}) async {
    try {
      await api.dio.post('/finance/iuran/mock', data: {
        'userId': userId,
        if (bulan != null) 'bulan': bulan,
        if (tahun != null) 'tahun': tahun,
      });
    } catch (e) {
      if (e is DioException) {
        throw Exception(e.response?.data['message'] ?? e.message ?? e.toString());
      }
      rethrow;
    }
  }

  Future<void> payIuran(String id, String method, String userId, {String? bulan, int? amount}) async {
    try {
      await api.dio.post('/finance/iuran/pay', data: {
        'id': id, 
        'method': method,
        'userId': userId,
        if (bulan != null) 'bulan': bulan,
        if (amount != null) 'amount': amount,
      });
    } catch (e) {
      if (e is DioException) {
        throw Exception(e.response?.data['message'] ?? e.message ?? e.toString());
      }
      rethrow;
    }
  }
}

class EventRepository {
  Future<List<Event>> getEvents() async {
    try {
      final response = await api.dio.get('/event');
      final List<dynamic> list = response.data['data'] ?? [];
      return list.map((item) => Event.fromJson(item)).toList();
    } catch (_) {
      return [
        Event(id: '1', jenis: 'Latgab', nama: 'Latihan Gabungan Nasional Jabar', lokasi: 'GOR Arcamanik Bandung', tanggal: '2026-07-28', deskripsi: 'Latihan gabungan bersama Guru Besar untuk tingkat Dasar s/d GPK.'),
        Event(id: '2', jenis: 'EKT', nama: 'EKT Jurus Yogyakarta', lokasi: 'Lapangan Kotagede', tanggal: '2026-08-15', deskripsi: 'Evaluasi Kenaikan Tingkat Jurus bagi segenap anggota Cabang Yogyakarta.'),
      ];
    }
  }

  Future<void> registerLatgab(Map<String, dynamic> data) async {
    await api.dio.post('/event/register', data: data);
  }

  Future<void> saveReservasi(Map<String, dynamic> data) async {
    await api.dio.post('/event/reservasi', data: data);
  }
}

class TransactionRepository {
  Future<void> saveNafas(Map<String, dynamic> data) async {
    await api.dio.post('/nafas/history', data: data);
  }

  Future<void> saveKebugaran(String userId, Map<String, dynamic> data) async {
    await api.dio.post('/organization/anggota/$userId/kebugaran', data: data);
  }

  Future<List<dynamic>> getKebugaranHistory(String userId) async {
    try {
      final response = await api.dio.get('/organization/anggota/$userId/kebugaran/history');
      return response.data['data'] ?? [];
    } catch (e) {
      return [];
    }
  }

  Future<void> saveAntrian(Map<String, dynamic> data) async {
    await api.dio.post('/training/antrian', data: data);
  }

  Future<Map<String, dynamic>> getAntrianStatus(String userId) async {
    final response = await api.dio.get('/training/antrian/status', queryParameters: {'user_id': userId});
    return response.data['data'];
  }

  Future<void> saveKta(Map<String, dynamic> data) async {
    await api.dio.post('/organization/anggota/kta', data: data);
  }

  Future<void> saveEkta(Map<String, dynamic> data) async {
    await api.dio.post('/organization/anggota/ekta', data: data);
  }

  Future<void> saveJadwal(Map<String, dynamic> data) async {
    await api.dio.post('/training/sesi', data: data);
  }
}
