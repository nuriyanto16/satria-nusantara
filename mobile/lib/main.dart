import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'core/theme.dart';
import 'data/repositories.dart';
import 'blocs/auth_bloc.dart';
import 'blocs/sesi_bloc.dart';
import 'blocs/iuran_bloc.dart';
import 'presentation/screens.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    // Instantiate repositories
    final authRepository = AuthRepository();
    final sesiRepository = SesiRepository();
    final financeRepository = FinanceRepository();

    return MultiBlocProvider(
      providers: [
        BlocProvider<AuthBloc>(
          create: (context) => AuthBloc(authRepository),
        ),
        BlocProvider<SesiBloc>(
          create: (context) => SesiBloc(sesiRepository),
        ),
        BlocProvider<IuranBloc>(
          create: (context) => IuranBloc(financeRepository),
        ),
      ],
      child: ValueListenableBuilder<ThemeMode>(
        valueListenable: themeNotifier,
        builder: (context, currentMode, child) {
          return MaterialApp(
            title: 'Satria Nusantara',
            debugShowCheckedModeBanner: false,
            theme: getBrandTheme(),
            darkTheme: getBrandDarkTheme(),
            themeMode: currentMode,
            initialRoute: '/splash',
            routes: {
              '/splash': (context) => const SplashScreen(),
              '/onboarding': (context) => const OnboardingScreen(),
              '/login': (context) => const LoginScreen(),
              '/home': (context) => const HomeScreen(),
              '/kta': (context) => const KtaDigitalScreen(),
              '/scan': (context) => const QrScannerScreen(),
              '/kebugaran': (context) => const KebugaranScoringScreen(),
              '/register': (context) => const RegisterWizardScreen(),
              '/event_detail': (context) => const EventDetailScreen(),
              '/event_participants': (context) => const EventParticipantListScreen(),
              '/google_complete': (context) => const GoogleDataCompleteScreen(),
              '/e_wallet_selection': (context) => const EWalletSelectionScreen(),
              '/payment_success': (context) => const PaymentSuccessScreen(),
              '/transfer_bukti': (context) => const TransferBuktiScreen(),
              '/sesi_detail': (context) => const SesiDetailScreen(),
              '/news_detail': (context) => const NewsDetailScreen(),
              '/kehadiran_detail': (context) => const KehadiranDetailScreen(),
              '/wait_verification': (context) => const WaitVerificationScreen(),
              '/nafas': (context) => const NafasDetailScreen(),
            },
          );
        },
      ),
    );
  }
}
