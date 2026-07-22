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
  }) async {
    final cleanEmail = email.trim().toLowerCase();
    final newUser = User(
      id: 'u-registered-${DateTime.now().millisecondsSinceEpoch}',
      email: cleanEmail,
      namaLengkap: name,
      noHp: phone,
      roleId: 4,
      roleName: 'Anggota',
      scope: 'anggota',
      status: 'pending', // Starts as pending verification!
    );
    _registeredUsers[cleanEmail] = newUser;
    _registeredPasswords[cleanEmail] = password;
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
  }

  Future<Map<String, dynamic>> loginGoogle(String email, String name, {String? googleId, String? noHp, String? fotoUrl}) async {
    try {
      final response = await api.dio.post(
        ApiConstants.googleLogin,
        data: {
          'email': email,
          'nama_lengkap': name,
          if (googleId != null && googleId.isNotEmpty) 'google_id': googleId,
          if (noHp != null && noHp.isNotEmpty) 'no_hp': noHp,
          if (fotoUrl != null && fotoUrl.isNotEmpty) 'foto_url': fotoUrl,
        },
      );
      final data = response.data['data'];
      final token = data['token'];
      final user = User.fromJson(data['user']);
      api.setToken(token);
      return {'token': token, 'user': user};
    } catch (e) {
      if (e is DioException && e.response?.statusCode == 404) {
        rethrow;
      }

      // Offline / Fallback handling for prototype testing:
      if (email != 'demo.anggota@gmail.com') {
        // If testing a custom email and noHp is not provided yet, force profile completion step
        if (noHp == null || noHp.isEmpty) {
          throw Exception("User not found");
        }
      }

      final mockToken = 'mock_google_token_${DateTime.now().millisecondsSinceEpoch}';
      final mockUser = User(
        id: 'u-google-${DateTime.now().millisecondsSinceEpoch}',
        namaLengkap: name,
        email: email,
        noHp: noHp ?? '081234567890',
        roleId: 4,
        roleName: 'Anggota',
        scope: 'anggota',
        status: 'aktif',
      );
      api.setToken(mockToken);
      return {'token': mockToken, 'user': mockUser};
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
    } catch (_) {
      // Fallback stub data for testing offline
      return [
        Sesi(
          id: '1',
          unitId: 'unit-01',
          unitNama: 'Unit Bandung',
          pelatihId: 'p001', pelatihNama: 'Pak Budi Susanto',
          jenis: 'Rutin',
          tanggal: '2026-07-20',
          jamMulai: '16:00',
          jamSelesai: '18:00',
          lokasi: 'Lapangan Gasibu Bandung',
          materi: 'Jurus Dasar 1-5',
        ),
        Sesi(
          id: '2',
          unitId: 'unit-01',
          unitNama: 'Unit Bandung',
          pelatihId: 'p001', pelatihNama: 'Pak Budi Susanto',
          jenis: 'Gabungan',
          tanggal: '2026-07-22',
          jamMulai: '15:30',
          jamSelesai: '17:30',
          lokasi: 'GOR Arcamanik Bandung',
          materi: 'Napas Pengolahan',
        ),
        Sesi(
          id: '3',
          unitId: 'unit-01',
          unitNama: 'Unit Bandung',
          pelatihId: 'p001', pelatihNama: 'Pak Budi Susanto',
          jenis: 'Rutin',
          tanggal: '2026-07-25',
          jamMulai: '16:00',
          jamSelesai: '18:00',
          lokasi: 'Lapangan Gasibu Bandung',
          materi: 'Jurus Praktis',
        ),
      ];
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
  static final Map<String, List<Iuran>> _userMockLists = {};

  Future<List<Iuran>> getIuranHistory(String userId) async {
    if (!_userMockLists.containsKey(userId)) {
      if (userId.contains('admin') || userId.contains('da001') || userId.contains('Sri') || userId.contains('Ahmad') || userId.contains('google')) {
        // Pre-existing members get full default history
        _userMockLists[userId] = [
          Iuran(id: '1', anggotaId: userId, bulan: 7, tahun: 2026, nominal: 40000, status: 'belum_bayar'),
          Iuran(id: '2', anggotaId: userId, bulan: 6, tahun: 2026, nominal: 40000, status: 'lunas', tanggalBayar: '2026-06-10'),
          Iuran(id: '3', anggotaId: userId, bulan: 5, tahun: 2026, nominal: 40000, status: 'lunas', tanggalBayar: '2026-05-12'),
        ];
      } else {
        // Custom registered user gets exactly 1 unpaid current month bill and no old history
        _userMockLists[userId] = [
          Iuran(id: '1', anggotaId: userId, bulan: 7, tahun: 2026, nominal: 40000, status: 'belum_bayar'),
        ];
      }
    }

    try {
      final response = await api.dio.get('/finance/iuran');
      final List<dynamic> list = response.data['data'] ?? [];
      final apiList = list.map((item) => Iuran.fromJson(item)).toList();
      if (apiList.isNotEmpty) {
        _userMockLists[userId] = apiList;
      }
      return _userMockLists[userId]!;
    } catch (_) {
      return _userMockLists[userId]!;
    }
  }

  Future<void> payIuran(String id, String method, String userId) async {
    try {
      await api.dio.post('/finance/iuran/pay', data: {
        'id': id, 
        'method': method,
        'userId': userId,
      });
    } catch (_) {}
    
    final list = _userMockLists[userId] ?? [];
    final index = list.indexWhere((item) => item.id == id);
    if (index != -1) {
      final old = list[index];
      list[index] = Iuran(
        id: old.id,
        anggotaId: old.anggotaId,
        bulan: old.bulan,
        tahun: old.tahun,
        nominal: old.nominal,
        status: 'lunas',
        tanggalBayar: DateTime.now().toString().split(' ')[0],
      );
    }
    await Future.delayed(const Duration(milliseconds: 800));
  }

  void addMockIuran(String userId) {
    int nextMonth = 8;
    int nextYear = 2026;
    final list = _userMockLists[userId] ?? [];
    if (list.isNotEmpty) {
      final last = list.first;
      nextMonth = last.bulan + 1;
      nextYear = last.tahun;
      if (nextMonth > 12) {
        nextMonth = 1;
        nextYear++;
      }
    }
    list.insert(0, Iuran(
      id: DateTime.now().millisecondsSinceEpoch.toString(),
      anggotaId: userId,
      bulan: nextMonth,
      tahun: nextYear,
      nominal: 40000,
      status: 'belum_bayar',
    ));
    _userMockLists[userId] = list;
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
