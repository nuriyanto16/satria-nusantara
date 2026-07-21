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

  static const Color bg = Color(0xFFF2F3F0);
  static const Color card = Colors.white;

  static const Color text1 = Color(0xFF1A1A18);
  static const Color text2 = Color(0xFF6B6B67);
  static const Color text3 = Color(0xFF9A9A96);
  static const Color border = Color(0x17000000);

  static const LinearGradient headerGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF0078C8), Color(0xFF0092EC), Color(0xFF42B8F5)],
  );
}

ThemeData getBrandTheme() {
  return ThemeData(
    primaryColor: BrandColors.biruPrimary,
    scaffoldBackgroundColor: BrandColors.bg,
    colorScheme: ColorScheme.fromSeed(
      seedColor: BrandColors.biruPrimary,
      primary: BrandColors.biruPrimary,
      secondary: BrandColors.hijau,
      background: BrandColors.bg,
      surface: BrandColors.card,
    ),
    useMaterial3: true,
    textTheme: GoogleFonts.poppinsTextTheme(),
    cardTheme: CardThemeData(
      color: BrandColors.card,
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
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
