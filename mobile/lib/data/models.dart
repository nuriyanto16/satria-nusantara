class User {
  final String id;
  final String email;
  final String namaLengkap;
  final String noHp;
  final int roleId;
  final String roleName;
  final String scope;
  final String status;

  User({
    required this.id,
    required this.email,
    required this.namaLengkap,
    required this.noHp,
    required this.roleId,
    required this.roleName,
    required this.scope,
    required this.status,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] ?? '',
      email: json['email'] ?? '',
      namaLengkap: json['nama_lengkap'] ?? '',
      noHp: json['no_hp'] ?? '',
      roleId: json['role_id'] ?? 0,
      roleName: json['role_name'] ?? '',
      scope: json['scope'] ?? '',
      status: json['status'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'nama_lengkap': namaLengkap,
      'no_hp': noHp,
      'role_id': roleId,
      'role_name': roleName,
      'scope': scope,
      'status': status,
    };
  }
}

class Sesi {
  final String id;
  final String unitId;
  final String unitNama;
  final String pelatihId;
  final String pelatihNama;
  final String jenis;
  final String tanggal;
  final String jamMulai;
  final String jamSelesai;
  final String lokasi;
  final String materi;

  Sesi({
    required this.id,
    required this.unitId,
    required this.unitNama,
    required this.pelatihId,
    required this.pelatihNama,
    required this.jenis,
    required this.tanggal,
    required this.jamMulai,
    required this.jamSelesai,
    required this.lokasi,
    required this.materi,
  });

  factory Sesi.fromJson(Map<String, dynamic> json) {
    String pNama = json['pelatih_nama'] ?? '';
    String pId = json['pelatih_id'] ?? '';
    if (pNama.isEmpty) {
      if (pId.isNotEmpty && !pId.contains('-')) {
        pNama = pId;
      } else {
        pNama = 'Pak Budi Susanto';
      }
    }
    return Sesi(
      id: json['id'] ?? '',
      unitId: json['unit_id'] ?? '',
      unitNama: json['unit_nama'] ?? '',
      pelatihId: pId,
      pelatihNama: pNama,
      jenis: json['jenis'] ?? '',
      tanggal: json['tanggal'] ?? '',
      jamMulai: json['jam_mulai'] ?? '',
      jamSelesai: json['jam_selesai'] ?? '',
      lokasi: json['lokasi'] ?? '',
      materi: json['materi'] ?? '',
    );
  }
}

class Iuran {
  final String id;
  final String anggotaId;
  final int bulan;
  final int tahun;
  final int nominal;
  final String status;
  final String? tanggalBayar;

  Iuran({
    required this.id,
    required this.anggotaId,
    required this.bulan,
    required this.tahun,
    required this.nominal,
    required this.status,
    this.tanggalBayar,
  });

  factory Iuran.fromJson(Map<String, dynamic> json) {
    return Iuran(
      id: json['id'] ?? '',
      anggotaId: json['anggota_id'] ?? '',
      bulan: json['bulan'] ?? 0,
      tahun: json['tahun'] ?? 0,
      nominal: json['nominal'] ?? 0,
      status: json['status'] ?? 'belum_bayar',
      tanggalBayar: json['tanggal_bayar'],
    );
  }
}

class Event {
  final String id;
  final String jenis; // Latgab, EKT, Pelatnas
  final String nama;
  final String lokasi;
  final String tanggal;
  final String deskripsi;

  Event({
    required this.id,
    required this.jenis,
    required this.nama,
    required this.lokasi,
    required this.tanggal,
    required this.deskripsi,
  });

  factory Event.fromJson(Map<String, dynamic> json) {
    return Event(
      id: json['id'] ?? '',
      jenis: json['jenis'] ?? '',
      nama: json['nama'] ?? '',
      lokasi: json['lokasi'] ?? '',
      tanggal: json['tanggal'] ?? '',
      deskripsi: json['deskripsi'] ?? '',
    );
  }
}
