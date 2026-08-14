import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:app_links/app_links.dart';
import 'core/theme.dart';
import 'data/repositories.dart';
import 'blocs/auth_bloc.dart';
import 'blocs/sesi_bloc.dart';
import 'blocs/iuran_bloc.dart';
import 'presentation/screens.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  final _navigatorKey = GlobalKey<NavigatorState>();
  late AppLinks _appLinks;

  @override
  void initState() {
    super.initState();
    _initDeepLinks();
  }

  void _initDeepLinks() {
    _appLinks = AppLinks();
    
    // Handle link when app is in foreground/background
    _appLinks.uriLinkStream.listen((uri) {
      debugPrint('Deep Link received: $uri');
      if (uri.host == 'payment-success') {
        _navigatorKey.currentState?.pushNamed('/payment_success');
      }
    });
  }

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
            navigatorKey: _navigatorKey,
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
              '/event_payment': (context) => const EventPaymentScreen(),
              '/nafas': (context) => const NafasDetailScreen(),
              '/reservasi': (context) => ReservasiScreen(),
              '/antrian': (context) => AntrianScreen(),
              '/payment_history': (context) => const PaymentHistoryScreen(),
            },
          );
        },
      ),
    );
  }
}
