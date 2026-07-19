import 'package:flutter_bloc/flutter_bloc.dart';
import '../data/models.dart';
import '../data/repositories.dart';

abstract class SesiEvent {}
class LoadSesiList extends SesiEvent {
  final String? unitId;
  final String? tanggal;
  final bool loadMore;
  LoadSesiList({this.unitId, this.tanggal, this.loadMore = false});
}
class ScanQREvent extends SesiEvent {
  final String code;
  final String userId;
  ScanQREvent(this.code, this.userId);
}

abstract class SesiState {}
class SesiInitial extends SesiState {}
class SesiLoading extends SesiState {}
class SesiLoaded extends SesiState {
  final List<Sesi> list;
  final bool hasMore;
  final int currentPage;
  SesiLoaded(this.list, {this.hasMore = true, this.currentPage = 1});
}
class SesiError extends SesiState {
  final String message;
  SesiError(this.message);
}

class QrScanSuccess extends SesiState {}
class QrScanError extends SesiState {
  final String message;
  QrScanError(this.message);
}

class SesiBloc extends Bloc<SesiEvent, SesiState> {
  final SesiRepository _repository;

  SesiBloc(this._repository) : super(SesiInitial()) {
    on<LoadSesiList>((event, emit) async {
      final isLoadMore = event.loadMore;
      List<Sesi> oldList = [];
      int pageToLoad = 1;
      
      if (isLoadMore && state is SesiLoaded) {
        oldList = (state as SesiLoaded).list;
        pageToLoad = (state as SesiLoaded).currentPage + 1;
      } else {
        emit(SesiLoading());
      }
      
      try {
        final list = await _repository.getSesiList(
          unitId: event.unitId,
          tanggal: event.tanggal,
        );
        
        final startIndex = (pageToLoad - 1) * 10;
        if (startIndex >= list.length) {
          emit(SesiLoaded(oldList, hasMore: false, currentPage: pageToLoad - 1));
          return;
        }
        
        final pageList = list.skip(startIndex).take(10).toList();
        final combinedList = [...oldList, ...pageList];
        final hasMore = combinedList.length < list.length;
        
        emit(SesiLoaded(combinedList, hasMore: hasMore, currentPage: pageToLoad));
      } catch (e) {
        emit(SesiError('Gagal memuat jadwal latihan.'));
      }
    });

    on<ScanQREvent>((event, emit) async {
      emit(SesiLoading());
      try {
        await _repository.scanQR(event.code, event.userId);
        emit(QrScanSuccess());
      } catch (e) {
        emit(QrScanError('Absensi gagal. Kode QR tidak valid atau kedaluwarsa.'));
      }
    });
  }
}
