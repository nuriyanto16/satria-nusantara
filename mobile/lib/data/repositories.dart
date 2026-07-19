import 'package:dio/dio.dart';
import '../core/network.dart';
import '../core/constants.dart';
import 'models.dart';

class AuthRepository {
  Future<Map<String, dynamic>> login(String email, String password) async {
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
  static List<Iuran> _mockList = [
    Iuran(id: '1', anggotaId: 'da001', bulan: 7, tahun: 2026, nominal: 50000, status: 'belum_bayar'),
    Iuran(id: '2', anggotaId: 'da001', bulan: 6, tahun: 2026, nominal: 50000, status: 'lunas', tanggalBayar: '2026-06-10'),
    Iuran(id: '3', anggotaId: 'da001', bulan: 5, tahun: 2026, nominal: 50000, status: 'lunas', tanggalBayar: '2026-05-12'),
  ];

  Future<List<Iuran>> getIuranHistory() async {
    try {
      final response = await api.dio.get('/finance/iuran');
      final List<dynamic> list = response.data['data'] ?? [];
      _mockList = list.map((item) => Iuran.fromJson(item)).toList();
      return _mockList;
    } catch (_) {
      return _mockList;
    }
  }

  Future<void> payIuran(String id, String method) async {
    try {
      await api.dio.post('/finance/iuran/pay', data: {'id': id, 'method': method});
    } catch (_) {}
    
    final index = _mockList.indexWhere((item) => item.id == id);
    if (index != -1) {
      final old = _mockList[index];
      _mockList[index] = Iuran(
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

  void addMockIuran() {
    int nextMonth = 8;
    int nextYear = 2026;
    if (_mockList.isNotEmpty) {
      final last = _mockList.first;
      nextMonth = last.bulan + 1;
      nextYear = last.tahun;
      if (nextMonth > 12) {
        nextMonth = 1;
        nextYear++;
      }
    }
    _mockList.insert(0, Iuran(
      id: DateTime.now().millisecondsSinceEpoch.toString(),
      anggotaId: 'da001',
      bulan: nextMonth,
      tahun: nextYear,
      nominal: 50000,
      status: 'belum_bayar',
    ));
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
