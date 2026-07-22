import 'package:flutter_bloc/flutter_bloc.dart';
import '../data/models.dart';
import '../data/repositories.dart';

abstract class AuthEvent {}
class LoginRequested extends AuthEvent {
  final String email;
  final String password;
  LoginRequested(this.email, this.password);
}
class LoggedIn extends AuthEvent {
  final String token;
  final User user;
  LoggedIn({required this.token, required this.user});
}
class LogoutRequested extends AuthEvent {}

abstract class AuthState {}
class AuthInitial extends AuthState {}
class AuthLoading extends AuthState {}
class Authenticated extends AuthState {
  final String token;
  final User user;
  Authenticated(this.token, this.user);
}
class Unauthenticated extends AuthState {
  final String? errorMessage;
  Unauthenticated({this.errorMessage});
}

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthRepository _repository;

  AuthBloc(this._repository) : super(AuthInitial()) {
    on<LoginRequested>((event, emit) async {
      emit(AuthLoading());
      try {
        final result = await _repository.login(event.email, event.password);
        emit(Authenticated(result['token'], result['user']));
      } catch (e) {
        final errorMsg = e.toString();
        if (errorMsg.contains("PENDING_VERIFICATION")) {
          final mockPendingUser = User(
            id: 'pending',
            email: event.email,
            namaLengkap: 'Calon Anggota',
            noHp: '-',
            roleId: 4,
            roleName: 'Anggota',
            scope: 'anggota',
            status: 'pending',
          );
          emit(Authenticated('pending_token', mockPendingUser));
        } else {
          emit(Unauthenticated(errorMessage: 'Email atau password salah.'));
        }
      }
    });

    on<LoggedIn>((event, emit) {
      emit(Authenticated(event.token, event.user));
    });

    on<LogoutRequested>((event, emit) {
      _repository.logout();
      emit(Unauthenticated());
    });
  }
}
