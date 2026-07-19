import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class BrandColors {
  static const Color hijau = Color(0xFF0F6B4C); // Premium Emerald Green
  static const Color hijauMuda = Color(0xFF1E8B65);
  static const Color hijauSoft = Color(0xFFE5F3EE);
  
  static const Color kuning = Color(0xFFE8C42A);
  static const Color kuningSoft = Color(0xFFFFF8E0);
  
  static const Color merah = Color(0xFFC0272D);
  static const Color biru = Color(0xFF5BB8D4);
  static const Color biruSoft = Color(0xFFE0F5FB);
  
  static const Color bg = Color(0xFFF2F3F0);
  static const Color card = Colors.white;
  
  static const Color text1 = Color(0xFF1A1A18);
  static const Color text2 = Color(0xFF6B6B67);
  static const Color text3 = Color(0xFF9A9A96);
  static const Color border = Color(0x17000000);
}

ThemeData getBrandTheme() {
  return ThemeData(
    primaryColor: BrandColors.hijau,
    scaffoldBackgroundColor: BrandColors.bg,
    colorScheme: ColorScheme.fromSeed(
      seedColor: BrandColors.hijau,
      primary: BrandColors.hijau,
      secondary: BrandColors.kuning,
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
      backgroundColor: BrandColors.hijau,
      foregroundColor: Colors.white,
      elevation: 0,
      centerTitle: true,
    ),
  );
}
