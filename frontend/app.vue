<template>
  <div class="layout">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<style>
/* Global CSS variables from mockup will go here or in main.css */
:root{
  --hijau:#1a5c2a;--hijau2:#2e8b3a;--hijau3:#e8f5ec;--hijau4:#a3d4b0;
  --kuning:#e8c42a;--merah:#c0272d;--biru:#5bb8d4;
  --bg:#f2f3f0;--card:#fff;--text1:#1a1a18;--text2:#6b6b67;--text3:#9a9a96;
  --border:rgba(0,0,0,0.09);--border2:rgba(0,0,0,0.14);--surface:#f5f5f2;
  --r6:6px;--r8:8px;--r12:12px;
}
html,body{height:100%;overflow:hidden;font-family:'Segoe UI',system-ui,sans-serif;background:var(--bg);color:var(--text1);margin:0;padding:0;box-sizing:border-box;}
.layout{display:flex;height:100vh;}
</style>
