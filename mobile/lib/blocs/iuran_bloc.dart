import 'package:flutter_bloc/flutter_bloc.dart';
import '../data/models.dart';
import '../data/repositories.dart';

abstract class IuranEvent {}
class LoadIuranHistory extends IuranEvent {}
class AddMockIuran extends IuranEvent {}
class PayIuranRequested extends IuranEvent {
  final String id;
  final String method;
  PayIuranRequested(this.id, this.method);
}

abstract class IuranState {}
class IuranInitial extends IuranState {}
class IuranLoading extends IuranState {}
class IuranLoaded extends IuranState {
  final List<Iuran> list;
  IuranLoaded(this.list);
}
class IuranError extends IuranState {
  final String message;
  IuranError(this.message);
}

class PaymentSuccess extends IuranState {}
class PaymentError extends IuranState {
  final String message;
  PaymentError(this.message);
}

class IuranBloc extends Bloc<IuranEvent, IuranState> {
  final FinanceRepository _repository;

  IuranBloc(this._repository) : super(IuranInitial()) {
    on<LoadIuranHistory>((event, emit) async {
      emit(IuranLoading());
      try {
        final list = await _repository.getIuranHistory();
        emit(IuranLoaded(list));
      } catch (e) {
        emit(IuranError('Gagal memuat histori iuran.'));
      }
    });

    on<AddMockIuran>((event, emit) async {
      emit(IuranLoading());
      try {
        _repository.addMockIuran();
        final list = await _repository.getIuranHistory();
        emit(IuranLoaded(list));
      } catch (e) {
        emit(IuranError('Gagal menambah tagihan.'));
      }
    });

    on<PayIuranRequested>((event, emit) async {
      emit(IuranLoading());
      try {
        await _repository.payIuran(event.id, event.method);
        emit(PaymentSuccess());
      } catch (e) {
        emit(PaymentError('Pembayaran gagal. Silakan coba lagi.'));
      }
    });
  }
}
