import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class BrandColors {
  // Primary Ocean Blue Theme matching SN_Mobile_Mockup_v10.html
  static const Color biruPrimary = Color(0xFF0092EC);
  static const Color biru = Color(0xFF0092EC);
  static const Color biruDark = Color(0xFF0078C8);
  static const Color biruLight = Color(0xFF42B8F5);
  static const Color biruSoft = Color(0xFFE8F6FF);

  // Map theme variables to Ocean Blue
  static const Color hijau = Color(0xFF0092EC);
  static const Color hijauMuda = Color(0xFF0078C8);
  static const Color hijauSoft = Color(0xFFE8F6FF);

  static const Color kuning = Color(0xFFE8C42A);
  static const Color kuningSoft = Color(0xFFFFF8E0);

  static const Color merah = Color(0xFFC0272D);

  // Light Mode Colors
  static const Color bg = Color(0xFFF2F3F0);
  static const Color card = Colors.white;
  static const Color text1 = Color(0xFF1A1A18);
  static const Color text2 = Color(0xFF6B6B67);
  static const Color text3 = Color(0xFF9A9A96);
  static const Color border = Color(0x17000000);

  // Dark Mode Colors
  static const Color bgDark = Color(0xFF121316);
  static const Color cardDark = Color(0xFF1E2026);
  static const Color surfaceDark = Color(0xFF262930);
  static const Color text1Dark = Color(0xFFF0F2F5);
  static const Color text2Dark = Color(0xFFA0A5B1);
  static const Color text3Dark = Color(0xFF6E7380);
  static const Color borderDark = Color(0x1FFFFFFF);

  static const LinearGradient headerGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF0078C8), Color(0xFF0092EC), Color(0xFF42B8F5)],
  );

  static const LinearGradient darkHeaderGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF0D253A), Color(0xFF143B5C), Color(0xFF0078C8)],
  );
}

// Global Theme Notifier for Dark/Light mode switching
class ThemeNotifier extends ValueNotifier<ThemeMode> {
  ThemeNotifier([ThemeMode mode = ThemeMode.light]) : super(mode);

  bool get isDarkMode => value == ThemeMode.dark;

  void toggleTheme() {
    value = value == ThemeMode.light ? ThemeMode.dark : ThemeMode.light;
  }

  void setTheme(ThemeMode mode) {
    value = mode;
  }
}

final themeNotifier = ThemeNotifier();

ThemeData getBrandTheme() {
  return ThemeData(
    brightness: Brightness.light,
    primaryColor: BrandColors.biruPrimary,
    scaffoldBackgroundColor: BrandColors.bg,
    colorScheme: ColorScheme.fromSeed(
      seedColor: BrandColors.biruPrimary,
      brightness: Brightness.light,
      primary: BrandColors.biruPrimary,
      secondary: BrandColors.hijau,
      background: BrandColors.bg,
      surface: BrandColors.card,
    ),
    useMaterial3: true,
    textTheme: GoogleFonts.poppinsTextTheme(ThemeData.light().textTheme),
    cardTheme: CardThemeData(
      color: BrandColors.card,
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
        side: const BorderSide(color: BrandColors.border, width: 1),
      ),
    ),
    appBarTheme: const AppBarTheme(
      backgroundColor: BrandColors.biruPrimary,
      foregroundColor: Colors.white,
      elevation: 0,
      centerTitle: true,
    ),
  );
}

ThemeData getBrandDarkTheme() {
  return ThemeData(
    brightness: Brightness.dark,
    primaryColor: BrandColors.biruPrimary,
    scaffoldBackgroundColor: BrandColors.bgDark,
    colorScheme: ColorScheme.fromSeed(
      seedColor: BrandColors.biruPrimary,
      brightness: Brightness.dark,
      primary: BrandColors.biruPrimary,
      secondary: BrandColors.biruLight,
      background: BrandColors.bgDark,
      surface: BrandColors.cardDark,
    ),
    useMaterial3: true,
    textTheme: GoogleFonts.poppinsTextTheme(ThemeData.dark().textTheme),
    cardTheme: CardThemeData(
      color: BrandColors.cardDark,
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
        side: const BorderSide(color: BrandColors.borderDark, width: 1),
      ),
    ),
    appBarTheme: const AppBarTheme(
      backgroundColor: BrandColors.cardDark,
      foregroundColor: Colors.white,
      elevation: 0,
      centerTitle: true,
    ),
  );
}
