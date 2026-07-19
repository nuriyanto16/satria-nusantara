import 'package:dio/dio.dart';
import 'constants.dart';

class NetworkClient {
  final Dio dio;
  String? _token;

  NetworkClient() : dio = Dio(BaseOptions(baseUrl: ApiConstants.baseUrl)) {
    dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) {
          options.baseUrl = ApiConstants.baseUrl; // Allow dynamic base URL
          if (_token != null && _token!.isNotEmpty) {
            options.headers['Authorization'] = 'Bearer $_token';
          }
          return handler.next(options);
        },
        onError: (DioException e, handler) {
          if (e.response?.statusCode == 401) {
            // Handle global logout or token refresh if necessary
          }
          return handler.next(e);
        },
      ),
    );
  }

  void setToken(String? token) {
    _token = token;
  }

  String? get token => _token;
}

// Global single instance of network client
final NetworkClient api = NetworkClient();
