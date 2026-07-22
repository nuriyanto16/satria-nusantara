import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class BrandColors {
  // Premium Deep Navy & Cyan Theme
  static const Color biruPrimary = Color(0xFF1E40AF); // Royal Sapphire Blue
  static const Color biru = Color(0xFF1E40AF);
  static const Color biruDark = Color(0xFF0F172A); // Midnight Slate
  static const Color biruLight = Color(0xFF00D2FF); // Radiant Energy Cyan
  static const Color biruSoft = Color(0xFFF0F9FF); // Clean Sky Blue

  // Map theme variables to Ocean Blue / Cyan
  static const Color hijau = Color(0xFF10B981); // Emerald Green
  static const Color hijauMuda = Color(0xFF059669);
  static const Color hijauSoft = Color(0xFFECFDF5);

  static const Color kuning = Color(0xFFF59E0B); // Amber Gold
  static const Color kuningSoft = Color(0xFFFEF3C7);

  static const Color merah = Color(0xFFEF4444); // Modern Coral Red

  // Light Mode Colors
  static const Color bg = Color(0xFFF8FAFC); // Clean Slate Gray
  static const Color card = Colors.white;
  static const Color text1 = Color(0xFF0F172A); // Midnight Slate text
  static const Color text2 = Color(0xFF475569); // Slate body text
  static const Color text3 = Color(0xFF94A3B8); // Gray placeholder text
  static const Color border = Color(0x0F000000);

  // Dark Mode Colors
  static const Color bgDark = Color(0xFF070B19); // Midnight deep background
  static const Color cardDark = Color(0xFF0F172A); // Slate dark cards
  static const Color surfaceDark = Color(0xFF1E293B);
  static const Color text1Dark = Color(0xFFF8FAFC);
  static const Color text2Dark = Color(0xFF94A3B8);
  static const Color text3Dark = Color(0xFF64748B);
  static const Color borderDark = Color(0x1FFFFFFF);

  static const LinearGradient headerGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF0F172A), Color(0xFF1E40AF), Color(0xFF00D2FF)],
  );

  static const LinearGradient darkHeaderGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF070B19), Color(0xFF0F1E3D), Color(0xFF1E40AF)],
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
