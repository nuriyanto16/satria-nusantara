import 'package:dio/dio.dart';
import '../core/network.dart';
import '../core/constants.dart';
import 'models.dart';

class AuthRepository {
  // Static map to hold newly registered mock users for prototype demo
  static final Map<String, User> _registeredUsers = {};
  static final Map<String, String> _registeredPasswords = {};

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
    
    // Pre-populate mock users list for offline/prototype fallback
    _registeredUsers[cleanEmail] = User(
      id: 'u-registered-${DateTime.now().millisecondsSinceEpoch}',
      email: cleanEmail,
      namaLengkap: name,
      noHp: phone,
      roleId: 4,
      roleName: 'Anggota',
      scope: 'anggota',
      status: 'pending',
    );
    _registeredPasswords[cleanEmail] = password;

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
    final cleanEmail = email.trim().toLowerCase();
    if (_registeredUsers.containsKey(cleanEmail)) {
      final user = _registeredUsers[cleanEmail]!;
      _registeredUsers[cleanEmail] = User(
        id: user.id,
        email: user.email,
        namaLengkap: user.namaLengkap,
        noHp: user.noHp,
        roleId: user.roleId,
        roleName: user.roleName,
        scope: user.scope,
        status: 'aktif',
      );
    }
  }

  Future<Map<String, dynamic>> login(String email, String password) async {
    final cleanEmail = email.trim().toLowerCase();

    // Check if it exists in mock registered users first
    if (_registeredUsers.containsKey(cleanEmail)) {
      if (_registeredPasswords[cleanEmail] == password) {
        final user = _registeredUsers[cleanEmail]!;
        final mockToken = 'mock_jwt_token_${DateTime.now().millisecondsSinceEpoch}';
        api.setToken(mockToken);
        return {'token': mockToken, 'user': user};
      } else {
        throw Exception("Email atau password salah");
      }
    }

    try {
      // Otherwise, fall back to real backend API call
      final response = await api.dio.post(
        ApiConstants.login,
        data: {'email': email, 'password': password},
      );
      final data = response.data['data'];
      final token = data['token'];
      final user = User.fromJson(data['user']);
      
      // Save token in API client
      api.setToken(token);
      
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

      // Offline / Fallback handling for prototype testing:
      // If the email is in mock memory store, return that user directly
      if (_registeredUsers.containsKey(cleanEmail)) {
        final user = _registeredUsers[cleanEmail]!;
        final mockToken = 'mock_google_token_${DateTime.now().millisecondsSinceEpoch}';
        api.setToken(mockToken);
        return {'token': mockToken, 'user': user};
      }

      // Otherwise, since this is a new Google login/registration, force profile completion
      throw Exception("User not found");
    }
  }

  void logout() {
    api.setToken(null);
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
}
