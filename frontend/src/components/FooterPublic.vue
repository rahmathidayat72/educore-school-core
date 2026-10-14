<template>
  <div>
    <div class="bg-dark text-white py-5">
      <div class="container">
        <div class="row g-4">
          <div class="col-lg-6 col-md-6">
            <div class="h5 text-uppercase fw-bold text-success mb-3">
              <i class="fas fa-graduation-cap me-2"></i>SD IT SEKOLAH KITA
            </div>
            <p class="text-secondary mb-3">
              Mendidik generasi berkarakter unggul, cerdas, berakhlak mulia, dan siap menghadapi tantangan masa depan dengan kurikulum berbasis teknologi dan nilai-nilai luhur.
            </p>
            <div class="d-flex gap-2">
              <a href="#" class="btn btn-sm btn-outline-light rounded-circle" style="width: 36px; height: 36px; display: inline-flex; align-items: center; justify-content: center;">
                <i class="fab fa-facebook-f"></i>
              </a>
              <a href="#" class="btn btn-sm btn-outline-light rounded-circle" style="width: 36px; height: 36px; display: inline-flex; align-items: center; justify-content: center;">
                <i class="fab fa-instagram"></i>
              </a>
              <a href="#" class="btn btn-sm btn-outline-light rounded-circle" style="width: 36px; height: 36px; display: inline-flex; align-items: center; justify-content: center;">
                <i class="fab fa-youtube"></i>
              </a>
            </div>
          </div>
          <div class="col-lg-3 col-md-6">
            <div class="h5 fw-bold mb-3">Navigasi Cepat</div>
            <ul class="list-unstyled">
              <li class="pb-2"><a href="#tentang" class="text-secondary text-decoration-none hover-white">Tentang Kami</a></li>
              <li class="pb-2"><a href="#keunggulan" class="text-secondary text-decoration-none hover-white">Fasilitas & Prestasi</a></li>
              <li class="pb-2"><a href="#guru" class="text-secondary text-decoration-none hover-white">Tenaga Pendidik</a></li>
              <li class="pb-2"><router-link to="/login" class="text-secondary text-decoration-none hover-white">Portal Akademik</router-link></li>
            </ul>
          </div>
          <div class="col-lg-3 col-md-6">
            <div class="h5 fw-bold mb-3">Kontak Sekolah</div>
            <p class="text-secondary mb-2 small"><i class="fas fa-map-marker-alt text-success me-2"></i>Jl. Pendidikan No. 123, Yogyakarta</p>
            <p class="text-secondary mb-2 small"><i class="fas fa-phone text-success me-2"></i>(0274) 123-4567</p>
            <p class="text-secondary mb-2 small"><i class="fab fa-whatsapp text-success me-2"></i>+62 812-3456-7890</p>
            <p class="text-secondary mb-2 small"><i class="fas fa-envelope text-success me-2"></i>info@sekolah.id</p>
          </div>
        </div>
      </div>
    </div>
    <div class="bg-black text-secondary py-3 border-top border-secondary small">
      <div class="container d-flex flex-column flex-md-row justify-content-between align-items-center gap-2">
        <div>Copyright &copy; 2026 Sistem Informasi Sekolah. All rights reserved.</div>
        <div>Terintegrasi dengan Go Native &amp; PostgreSQL</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.hover-white:hover {
  color: #ffffff !important;
}
</style>
