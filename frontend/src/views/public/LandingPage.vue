<template>
  <div class="landing-page">
    <NavbarPublic />

    <!-- Hero Carousel dengan Animasi Ken Burns & Staggered Captions -->
    <div id="heroCarousel" class="carousel slide carousel-fade" data-bs-ride="carousel" data-bs-interval="6500">
      <div class="carousel-indicators">
        <button type="button" data-bs-target="#heroCarousel" data-bs-slide-to="0" class="active" aria-current="true" aria-label="Slide 1"></button>
        <button type="button" data-bs-target="#heroCarousel" data-bs-slide-to="1" aria-label="Slide 2"></button>
        <button type="button" data-bs-target="#heroCarousel" data-bs-slide-to="2" aria-label="Slide 3"></button>
      </div>
      <div class="carousel-inner overflow-hidden">
        <!-- Slide 1 -->
        <div class="carousel-item active position-relative hero-slide">
          <div class="hero-img-container">
            <img
              src="https://naevaschool.naevaweb.my.id/userfiles/uploads/bg-edukasi-45.jpg"
              class="d-block w-100 object-fit-cover hero-kenburns-img"
              alt="Banner 1"
            />
          </div>
          <div class="hero-overlay"></div>
          <div class="carousel-caption text-start pb-5 hero-caption">
            <span class="badge bg-success mb-3 px-3 py-2 fs-6 rounded-pill shadow-sm animate-badge">
              <i class="fas fa-graduation-cap me-1"></i>Tahun Ajaran 2026/2027
            </span>
            <h1 class="display-4 fw-bold text-white mb-2 hero-title">Selamat Datang di Sekolah Unggulan</h1>
            <p class="lead text-light mb-4 hero-lead">Membentuk generasi pemimpin beriman, berilmu, dan berdaya saing global.</p>
            <div class="d-flex flex-wrap gap-3 hero-btns">
              <a href="#tentang" class="btn btn-success btn-lg px-4 py-2 shadow-sm btn-living">Pelajari Lebih Lanjut</a>
              <router-link to="/login" class="btn btn-outline-light btn-lg px-4 py-2 btn-living">
                <i class="fas fa-sign-in-alt me-1"></i>Masuk Portal
              </router-link>
            </div>
          </div>
        </div>

        <!-- Slide 2 -->
        <div class="carousel-item position-relative hero-slide">
          <div class="hero-img-container">
            <img
              src="https://naevaschool.naevaweb.my.id/userfiles/uploads/bg-edukasi-76.jpg"
              class="d-block w-100 object-fit-cover hero-kenburns-img hero-kenburns-alt"
              alt="Banner 2"
            />
          </div>
          <div class="hero-overlay"></div>
          <div class="carousel-caption text-start pb-5 hero-caption">
            <span class="badge bg-primary mb-3 px-3 py-2 fs-6 rounded-pill shadow-sm animate-badge">
              <i class="fas fa-award me-1"></i>Fasilitas Lengkap
            </span>
            <h1 class="display-4 fw-bold text-white mb-2 hero-title">Pembelajaran Modern &amp; Interaktif</h1>
            <p class="lead text-light mb-4 hero-lead">Didukung laboratorium modern, ruang multimedia, dan lingkungan asri ramah anak.</p>
            <div class="d-flex flex-wrap gap-3 hero-btns">
              <a href="#keunggulan" class="btn btn-primary btn-lg px-4 py-2 shadow-sm btn-living">Lihat Fasilitas</a>
              <router-link to="/login" class="btn btn-outline-light btn-lg px-4 py-2 btn-living">Portal Guru &amp; Siswa</router-link>
            </div>
          </div>
        </div>

        <!-- Slide 3 -->
        <div class="carousel-item position-relative hero-slide">
          <div class="hero-img-container">
            <img
              src="https://naevaschool.naevaweb.my.id/userfiles/uploads/bg-edukasi-86.jpg"
              class="d-block w-100 object-fit-cover hero-kenburns-img"
              alt="Banner 3"
            />
          </div>
          <div class="hero-overlay"></div>
          <div class="carousel-caption text-start pb-5 hero-caption">
            <span class="badge bg-warning text-dark mb-3 px-3 py-2 fs-6 rounded-pill shadow-sm animate-badge">
              <i class="fas fa-star me-1"></i>Ekstrakurikuler Aktif
            </span>
            <h1 class="display-4 fw-bold text-white mb-2 hero-title">Wadah Bakat dan Prestasi Siswa</h1>
            <p class="lead text-light mb-4 hero-lead">Mengembangkan minat siswa di bidang sains, olahraga, seni, dan kepemimpinan.</p>
            <div class="d-flex flex-wrap gap-3 hero-btns">
              <a href="#kontak" class="btn btn-warning text-dark btn-lg px-4 py-2 shadow-sm btn-living">Hubungi Kami</a>
              <a href="#tentang" class="btn btn-outline-light btn-lg px-4 py-2 btn-living">Kenali Sekolah</a>
            </div>
          </div>
        </div>
      </div>

      <button class="carousel-control-prev" type="button" data-bs-target="#heroCarousel" data-bs-slide="prev">
        <span class="carousel-control-prev-icon rounded-circle p-3 bg-dark bg-opacity-25" aria-hidden="true"></span>
        <span class="visually-hidden">Previous</span>
      </button>
      <button class="carousel-control-next" type="button" data-bs-target="#heroCarousel" data-bs-slide="next">
        <span class="carousel-control-next-icon rounded-circle p-3 bg-dark bg-opacity-25" aria-hidden="true"></span>
        <span class="visually-hidden">Next</span>
      </button>
    </div>

    <!-- Sambutan & Profil Singkat dengan Gambar Hidup -->
    <section id="tentang" class="container py-5">
      <div class="row py-3 align-items-center g-5">
        <div class="col-lg-7 reveal-on-scroll reveal-fade-left">
          <span class="badge bg-success-subtle text-success fw-bold text-uppercase px-3 py-1 mb-2">Profil Lembaga</span>
          <h2 class="display-6 fw-bold text-dark mt-1">SD IT Naeva School</h2>
          <h5 class="text-secondary mb-4 fw-medium">Menjadi Sekolah yang Ramah, Nyaman, dan Menyenangkan bagi Siswa</h5>
          <p class="text-muted leading-relaxed fs-6">
            Kami berkomitmen memberikan pendidikan dasar terbaik dengan memadukan kurikulum nasional dan pembentukan karakter akhlakul karimah. Didukung tenaga pendidik berpengalaman, sarana dan prasarana memadai, serta atmosfer belajar yang inspiratif.
          </p>
          <div class="d-flex flex-wrap gap-3 pt-3">
            <a href="#kontak" class="btn btn-dark px-4 py-2 btn-living shadow-sm">
              <i class="fas fa-info-circle me-1"></i>Informasi Pendaftaran
            </a>
            <router-link to="/login" class="btn btn-outline-dark px-4 py-2 btn-living">
              <i class="fas fa-school me-1"></i>Portal Akademik
            </router-link>
          </div>
        </div>

        <!-- Gambar Gedung Hidup dengan Floating Badge -->
        <div class="col-lg-5 text-center reveal-on-scroll reveal-fade-right">
          <div class="school-img-wrapper position-relative mx-auto">
            <!-- Glow Dekoratif di Belakang -->
            <div class="glow-orb"></div>
            
            <div class="img-living-wrapper shadow-lg rounded-4 overflow-hidden animate-float">
              <img
                src="https://naevaschool.naevaweb.my.id/userfiles/uploads/bg-edukasi-246.jpg"
                class="d-block w-100 object-fit-cover living-photo"
                alt="Gedung Sekolah"
                style="height: 380px;"
              />
            </div>

            <!-- Floating Badge Terakreditasi A -->
            <div class="floating-stat-card shadow-lg p-3 bg-white rounded-3 position-absolute d-flex align-items-center gap-3">
              <div class="bg-success text-white rounded-circle p-2 d-flex align-items-center justify-content-center" style="width: 44px; height: 44px;">
                <i class="fas fa-check-double fs-5"></i>
              </div>
              <div class="text-start">
                <div class="fw-bold text-dark small">Akreditasi A Unggul</div>
                <div class="text-secondary small" style="font-size: 0.75rem;">Standar Mutu Nasional</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Keunggulan Sekolah (Card 4 Kolom Hidup & Interaktif) -->
    <section id="keunggulan" class="bg-light py-5 position-relative overflow-hidden">
      <div class="container py-3">
        <div class="text-center mb-5 reveal-on-scroll">
          <span class="badge bg-success-subtle text-success fw-bold text-uppercase px-3 py-1 mb-2">Mengapa Memilih Kami?</span>
          <h2 class="fw-bold text-dark mt-1">Keunggulan &amp; Fasilitas Sekolah</h2>
          <p class="text-secondary max-w-700 mx-auto">Standar mutu terbaik untuk masa depan gemilang buah hati Anda dengan lingkungan pendukung yang inspiratif.</p>
        </div>

        <div class="row g-4">
          <!-- Card 1 -->
          <div class="col-lg-3 col-md-6 reveal-on-scroll delay-100">
            <div class="bg-white text-center p-4 rounded-4 border h-100 shadow-sm card-living-shine feature-card">
              <div class="icon-bubble text-success pb-3 mx-auto">
                <i class="fas fa-trophy fa-3x icon-living-trophy"></i>
              </div>
              <div class="h5 fw-bold text-uppercase pb-2 text-dark">Berprestasi</div>
              <p class="text-muted small mb-0">Ratusan prestasi akademik dan non-akademik di tingkat kota, provinsi, hingga nasional.</p>
            </div>
          </div>

          <!-- Card 2 -->
          <div class="col-lg-3 col-md-6 reveal-on-scroll delay-200">
            <div class="bg-white text-center p-4 rounded-4 border h-100 shadow-sm card-living-shine feature-card">
              <div class="icon-bubble text-primary pb-3 mx-auto">
                <i class="fas fa-award fa-3x icon-living-award"></i>
              </div>
              <div class="h5 fw-bold text-uppercase pb-2 text-dark">Terakreditasi A</div>
              <p class="text-muted small mb-0">Pengakuan mutu pendidikan dengan predikat akreditasi terbaik secara berkelanjutan.</p>
            </div>
          </div>

          <!-- Card 3 -->
          <div class="col-lg-3 col-md-6 reveal-on-scroll delay-300">
            <div class="bg-white text-center p-4 rounded-4 border h-100 shadow-sm card-living-shine feature-card">
              <div class="icon-bubble text-info pb-3 mx-auto">
                <i class="fas fa-microscope fa-3x icon-living-microscope"></i>
              </div>
              <div class="h5 fw-bold text-uppercase pb-2 text-dark">Fasilitas Modern</div>
              <p class="text-muted small mb-0">Laboratorium komputer, sains terpadu, perpustakaan digital, serta sarana olahraga lengkap.</p>
            </div>
          </div>

          <!-- Card 4 -->
          <div class="col-lg-3 col-md-6 reveal-on-scroll delay-400">
            <div class="bg-white text-center p-4 rounded-4 border h-100 shadow-sm card-living-shine feature-card">
              <div class="icon-bubble text-warning pb-3 mx-auto">
                <i class="fas fa-shapes fa-3x icon-living-shapes"></i>
              </div>
              <div class="h5 fw-bold text-uppercase pb-2 text-dark">Ekstrakurikuler</div>
              <p class="text-muted small mb-0">Pengembangan bakat holistik mulai dari robotik, pramuka, tahfidz, seni musik, hingga bela diri.</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Galeri & Artikel Kegiatan (Gambar Hidup) -->
    <section class="container py-5">
      <div class="row align-items-stretch g-4">
        <!-- Galeri Slider Hidup -->
        <div class="col-lg-5 reveal-on-scroll reveal-fade-left">
          <div class="bg-dark text-white p-4 rounded-4 h-100 d-flex flex-column justify-content-between shadow gallery-card position-relative overflow-hidden">
            <div>
              <div class="d-flex align-items-center justify-content-between pb-2">
                <div class="h5 text-uppercase fw-bold text-success mb-0">Galeri Aktivitas</div>
                <span class="badge bg-success-subtle text-success px-2 py-1 small">Terbaru</span>
              </div>
              <p class="small text-secondary mb-3">Momen belajar-mengajar dan keseruan siswa di sekolah.</p>
              
              <div id="carouselGallery" class="carousel slide carousel-fade rounded-3 overflow-hidden shadow-sm" data-bs-ride="carousel">
                <div class="carousel-inner">
                  <div class="carousel-item active">
                    <div class="img-living-wrapper">
                      <img src="https://naevaschool.naevaweb.my.id/userfiles/uploads/img-edukasi-64.jpg" class="d-block w-100 object-fit-cover gallery-img" style="height: 270px;" alt="Galeri 1" />
                    </div>
                  </div>
                  <div class="carousel-item">
                    <div class="img-living-wrapper">
                      <img src="https://naevaschool.naevaweb.my.id/userfiles/uploads/img-edukasi-5-4.jpg" class="d-block w-100 object-fit-cover gallery-img" style="height: 270px;" alt="Galeri 2" />
                    </div>
                  </div>
                  <div class="carousel-item">
                    <div class="img-living-wrapper">
                      <img src="https://naevaschool.naevaweb.my.id/userfiles/uploads/img-edukasi-57.jpg" class="d-block w-100 object-fit-cover gallery-img" style="height: 270px;" alt="Galeri 3" />
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="pt-3 text-center">
              <small class="text-secondary"><i class="fas fa-camera me-1"></i>Dokumentasi kegiatan siswa tahun ajaran aktif</small>
            </div>
          </div>
        </div>

        <!-- Berita & Artikel dengan Thumbnail Hidup -->
        <div class="col-lg-7 reveal-on-scroll reveal-fade-right">
          <div class="h5 text-uppercase fw-bold text-dark pb-2 border-bottom d-flex align-items-center justify-content-between">
            <span>Berita &amp; Informasi Sekolah</span>
            <span class="text-secondary small fw-normal"><i class="far fa-newspaper me-1"></i>Update Terkini</span>
          </div>

          <!-- Berita 1 -->
          <div class="news-item d-flex py-3 border-bottom align-items-center">
            <div class="img-living-wrapper rounded-3 me-3 flex-shrink-0" style="width: 82px; height: 82px;">
              <img src="https://naevaschool.naevaweb.my.id/userfiles/uploads/img-edukasi-57.jpg" class="w-100 h-100 object-fit-cover" alt="Berita 1" />
            </div>
            <div>
              <a href="#" class="text-decoration-none text-dark fw-bold news-link">Jadwal Pendaftaran dan Observasi Calon Siswa Baru (PPDB)</a>
              <p class="small text-secondary mb-0 mt-1">Pendaftaran gelombang pertama dibuka mulai bulan ini. Segera konfirmasi kuota kelas.</p>
              <div class="small text-success mt-1 fw-medium"><i class="far fa-clock me-1"></i>Oktober 2026</div>
            </div>
          </div>

          <!-- Berita 2 -->
          <div class="news-item d-flex py-3 border-bottom align-items-center">
            <div class="img-living-wrapper rounded-3 me-3 flex-shrink-0" style="width: 82px; height: 82px;">
              <img src="https://naevaschool.naevaweb.my.id/userfiles/uploads/img-edukasi-64.jpg" class="w-100 h-100 object-fit-cover" alt="Berita 2" />
            </div>
            <div>
              <a href="#" class="text-decoration-none text-dark fw-bold news-link">Pertemuan Wali Murid &amp; Workshop Parenting Digital</a>
              <p class="small text-secondary mb-0 mt-1">Kolaborasi aktif orang tua dan guru dalam mendampingi tumbuh kembang anak di era digital.</p>
              <div class="small text-success mt-1 fw-medium"><i class="far fa-clock me-1"></i>September 2026</div>
            </div>
          </div>

          <!-- Berita 3 -->
          <div class="news-item d-flex py-3 align-items-center">
            <div class="img-living-wrapper rounded-3 me-3 flex-shrink-0" style="width: 82px; height: 82px;">
              <img src="https://naevaschool.naevaweb.my.id/userfiles/uploads/img-edukasi-5-4.jpg" class="w-100 h-100 object-fit-cover" alt="Berita 3" />
            </div>
            <div>
              <a href="#" class="text-decoration-none text-dark fw-bold news-link">Pekan Kreativitas Seni &amp; Sains Siswa</a>
              <p class="small text-secondary mb-0 mt-1">Pameran karya inovatif siswa kelas 1 hingga kelas 6 dengan antusiasme yang luar biasa.</p>
              <div class="small text-success mt-1 fw-medium"><i class="far fa-clock me-1"></i>Agustus 2026</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Banner PPDB Call To Action Hidup & Berdenyut Halus -->
    <section class="py-4 reveal-on-scroll reveal-zoom">
      <div class="container">
        <div class="cta-banner text-white p-5 rounded-4 text-center shadow-lg position-relative overflow-hidden">
          <div class="cta-shine"></div>
          <span class="badge bg-white text-success px-3 py-2 fw-bold text-uppercase mb-3 pulse-beacon">
            <i class="fas fa-sparkles me-1"></i>Penerimaan Siswa Baru
          </span>
          <h2 class="display-6 fw-bold mb-3">Segera Daftarkan Putra-Putri Anda!</h2>
          <p class="lead mb-4 max-w-700 mx-auto text-light">
            Kuota terbatas untuk kelas baru tahun ajaran 2026/2027. Dapatkan beasiswa prestasi dan potongan biaya pendaftaran awal.
          </p>
          <div class="d-flex justify-content-center flex-wrap gap-3">
            <a href="#kontak" class="btn btn-light text-success fw-bold px-4 py-2 btn-living shadow-sm">
              <i class="fas fa-phone-alt me-1"></i>Info Lebih Lanjut
            </a>
            <router-link to="/login" class="btn btn-outline-light px-4 py-2 btn-living">
              <i class="fas fa-arrow-right me-1"></i>Akses Portal
            </router-link>
          </div>
        </div>
      </div>
    </section>

    <!-- Peta & Kontak Kami -->
    <section id="kontak" class="container py-5">
      <div class="row g-4 align-items-center">
        <div class="col-lg-6 reveal-on-scroll reveal-fade-left">
          <div class="rounded-4 overflow-hidden shadow-sm border map-wrapper">
            <iframe
              src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d63245.983719423304!2d110.33364512110442!3d-7.803163419271912!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e7a5787bd5b6bc5%3A0x21723fd4d3684f71!2sYogyakarta%2C%20Yogyakarta%20City%2C%20Special%20Region%20of%20Yogyakarta!5e0!3m2!1sen!2sid!4v1710005466701!5m2!1sen!2sid"
              width="100%"
              height="350"
              style="border:0;"
              allowfullscreen=""
              loading="lazy"
              referrerpolicy="no-referrer-when-downgrade"
            ></iframe>
          </div>
        </div>
        <div class="col-lg-6 reveal-on-scroll reveal-fade-right">
          <span class="badge bg-success-subtle text-success fw-bold text-uppercase px-3 py-1 mb-2">Hubungi Kami</span>
          <h2 class="fw-bold text-dark mt-1 mb-3">Layanan Informasi &amp; Kunjungan</h2>
          <p class="text-muted mb-4">Silakan hubungi kami atau berkunjung langsung ke kampus sekolah kami pada jam kerja (Senin - Jumat 07.30 - 15.30 WIB).</p>

          <div class="contact-pill d-flex align-items-center mb-3 p-2 rounded-3 border bg-white shadow-sm">
            <div class="bg-success text-white rounded-circle p-2 me-3 flex-shrink-0 d-flex align-items-center justify-content-center" style="width: 44px; height: 44px;">
              <i class="fas fa-map-marker-alt fs-5"></i>
            </div>
            <div>
              <div class="fw-bold text-dark small">Alamat Kampus</div>
              <div class="text-secondary small">Jl. Pendidikan No. 123, Sleman, D.I. Yogyakarta</div>
            </div>
          </div>

          <div class="contact-pill d-flex align-items-center mb-3 p-2 rounded-3 border bg-white shadow-sm">
            <div class="bg-success text-white rounded-circle p-2 me-3 flex-shrink-0 d-flex align-items-center justify-content-center" style="width: 44px; height: 44px;">
              <i class="fas fa-phone fs-5"></i>
            </div>
            <div>
              <div class="fw-bold text-dark small">Telepon Kantor</div>
              <div class="text-secondary small">(0274) 123-4567</div>
            </div>
          </div>

          <div class="contact-pill d-flex align-items-center mb-3 p-2 rounded-3 border bg-white shadow-sm">
            <div class="bg-success text-white rounded-circle p-2 me-3 flex-shrink-0 d-flex align-items-center justify-content-center" style="width: 44px; height: 44px;">
              <i class="fab fa-whatsapp fs-5"></i>
            </div>
            <div>
              <div class="fw-bold text-dark small">WhatsApp Helpdesk</div>
              <div class="text-secondary small">+62 812-3456-7890</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <FooterPublic />
  </div>
</template>

<script setup>
import NavbarPublic from '@/components/NavbarPublic.vue'
import FooterPublic from '@/components/FooterPublic.vue'
import { useScrollReveal } from '@/composables/useScrollReveal'

// Aktifkan scroll reveal untuk elemen yang masuk layar
useScrollReveal('.reveal-on-scroll', { threshold: 0.15 })
</script>

<style scoped>
/* =========================================
   HERO CAROUSEL & LIVING KEN BURNS EFFECT
   ========================================= */
.hero-slide {
  max-height: 520px;
  overflow: hidden;
}

.hero-img-container {
  width: 100%;
  height: 520px;
  overflow: hidden;
}

.hero-kenburns-img {
  height: 520px;
  filter: brightness(0.62) contrast(1.05);
  animation: heroKenBurns 18s ease-in-out infinite alternate;
  will-change: transform;
}

.hero-kenburns-alt {
  animation: heroKenBurnsAlt 20s ease-in-out infinite alternate;
}

@keyframes heroKenBurns {
  0% {
    transform: scale(1) translate(0, 0);
  }
  50% {
    transform: scale(1.09) translate(-1%, -1.2%);
  }
  100% {
    transform: scale(1.15) translate(1.2%, 0.8%);
  }
}

@keyframes heroKenBurnsAlt {
  0% {
    transform: scale(1.14) translate(1%, 0.5%);
  }
  50% {
    transform: scale(1.06) translate(-0.5%, -0.8%);
  }
  100% {
    transform: scale(1) translate(-1.2%, 0);
  }
}

.hero-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(180deg, rgba(15, 23, 42, 0.25) 0%, rgba(15, 23, 42, 0.65) 100%);
  pointer-events: none;
}

.hero-caption {
  z-index: 10;
  bottom: 10%;
}

.hero-title {
  animation: captionFadeUp 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.hero-lead {
  animation: captionFadeUp 1s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.hero-btns {
  animation: captionFadeUp 1.2s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes captionFadeUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* =========================================
   FOTO GEDUNG SEKOLAH HIDUP
   ========================================= */
.school-img-wrapper {
  position: relative;
  max-width: 440px;
}

.glow-orb {
  position: absolute;
  width: 320px;
  height: 320px;
  background: radial-gradient(circle, rgba(15, 118, 110, 0.22) 0%, rgba(15, 118, 110, 0) 70%);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 0;
  pointer-events: none;
  animation: pulseGlow 4s ease-in-out infinite alternate;
}

@keyframes pulseGlow {
  0% { transform: translate(-50%, -50%) scale(0.9); opacity: 0.6; }
  100% { transform: translate(-50%, -50%) scale(1.15); opacity: 0.95; }
}

.living-photo {
  transition: transform 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

.img-living-wrapper:hover .living-photo {
  transform: scale(1.07);
}

.floating-stat-card {
  bottom: -20px;
  left: -20px;
  z-index: 2;
  border: 1px solid #e2e8f0;
  animation: floatGentle 4.5s ease-in-out infinite 0.5s;
}

/* =========================================
   FEATURE CARDS & LIVING ICONS
   ========================================= */
.feature-card {
  transition: transform 0.35s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.35s ease, border-color 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 16px 32px rgba(0, 0, 0, 0.08) !important;
  border-color: #cbd5e1 !important;
}

.icon-bubble {
  transition: transform 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.feature-card:hover .icon-bubble {
  transform: scale(1.12);
}

.icon-living-trophy {
  animation: iconTrophyWiggle 5s ease-in-out infinite;
}

@keyframes iconTrophyWiggle {
  0%, 100% { transform: rotate(0deg); }
  25% { transform: rotate(-6deg); }
  75% { transform: rotate(6deg); }
}

.icon-living-award {
  animation: floatGentle 4s ease-in-out infinite;
}

.icon-living-microscope {
  animation: iconMicroscope 6s ease-in-out infinite;
}

@keyframes iconMicroscope {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.08) rotate(4deg); }
}

.icon-living-shapes {
  animation: breathSubtle 4s ease-in-out infinite;
}

/* =========================================
   NEWS ARTICLES & INTERACTIVE LINKS
   ========================================= */
.news-item {
  transition: all 0.25s ease;
  border-radius: 8px;
  padding-left: 6px;
  padding-right: 6px;
}

.news-item:hover {
  background-color: #f8fafc;
  transform: translateX(6px);
}

.news-link {
  transition: color 0.2s ease;
}

.news-item:hover .news-link {
  color: #198754 !important;
}

/* =========================================
   CTA BANNER
   ========================================= */
.cta-banner {
  background: linear-gradient(135deg, #0f766e 0%, #15803d 100%);
  position: relative;
}

.cta-banner::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.12) 0%, transparent 60%);
  pointer-events: none;
  animation: rotateGlow 14s linear infinite;
}

@keyframes rotateGlow {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* =========================================
   CONTACT PILLS
   ========================================= */
.contact-pill {
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.25s ease;
}

.contact-pill:hover {
  transform: translateX(6px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.06) !important;
}

.max-w-700 {
  max-width: 700px;
}
</style>
