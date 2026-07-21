class ApiConstants {
  static const String devBaseUrl = 'http://localhost:8080/api/v1';
  static const String androidEmulatorBaseUrl = 'http://10.0.2.2:8080/api/v1';
  static const String vpsBaseUrl = 'https://nfmtech.my.id/product/satrianusantara/api/v1';

  // Current active base URL
  static String baseUrl = vpsBaseUrl;

  static const String login = '/auth/login';
  static const String googleLogin = '/auth/google-login';
  static const String getSesi = '/training/sesi';
  static const String scanQr = '/training/scan-qr';
  static const String generateQr = '/training/sesi'; // dynamic path
}
