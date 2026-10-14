# Panduan Penggunaan `AGENTS.md`

Repository ini menyediakan [`AGENTS.md`](./AGENTS.md), sebuah standar kerja universal untuk membantu AI coding agent dan developer mengelola proyek perangkat lunak secara terstruktur—mulai dari ide, discovery, perencanaan, implementasi, verifikasi, release, hingga maintenance.

`AGENTS.md` bukan source code aplikasi dan bukan template teknologi tertentu. File tersebut berfungsi sebagai **protokol kerja permanen** yang mengatur bagaimana AI memahami proyek, meminta persetujuan, mengubah file, membuktikan hasil, dan memelihara dokumentasi proyek.

> **Sumber normatif:** README ini adalah panduan ringkas. Jika terdapat perbedaan atau ringkasan yang tidak lengkap, ikuti `AGENTS.md` lalu sinkronkan README.

> **Ringkasnya:** letakkan `AGENTS.md` di root proyek, buka AI coding agent dari root tersebut, lalu minta AI mengaudit dan menjalankan alur yang sesuai dengan kondisi proyek.

---

## Daftar Isi

- [Tujuan](#tujuan)
- [Konsep Utama](#konsep-utama)
- [Kapan Digunakan](#kapan-digunakan)
- [Persiapan](#persiapan)
- [Quick Start](#quick-start)
- [Penggunaan untuk Proyek Baru](#penggunaan-untuk-proyek-baru)
- [Penggunaan untuk Existing Repository](#penggunaan-untuk-existing-repository)
- [Melanjutkan Pekerjaan pada Sesi Berikutnya](#melanjutkan-pekerjaan-pada-sesi-berikutnya)
- [Tiga Gerbang Persetujuan](#tiga-gerbang-persetujuan)
- [Model Keamanan Menyeluruh](#model-keamanan-menyeluruh)
- [Dokumentasi yang Dihasilkan](#dokumentasi-yang-dihasilkan)
- [Siklus Pengerjaan Task](#siklus-pengerjaan-task)
- [Prompt Siap Pakai](#prompt-siap-pakai)
- [Bug, Perubahan Scope, dan Incident](#bug-perubahan-scope-dan-incident)
- [Release dan Deployment](#release-dan-deployment)
- [Tanggung Jawab Pengguna dan AI](#tanggung-jawab-pengguna-dan-ai)
- [Status dan Tracking](#status-dan-tracking)
- [Praktik yang Direkomendasikan](#praktik-yang-direkomendasikan)
- [Hal yang Dilarang atau Perlu Dihindari](#hal-yang-dilarang-atau-perlu-dihindari)
- [FAQ](#faq)
- [Troubleshooting](#troubleshooting)
- [Quick Reference](#quick-reference)

---

## Tujuan

`AGENTS.md` dirancang untuk mencegah pola kerja AI yang berisiko, seperti:

- langsung membuat aplikasi sebelum masalah bisnis dipahami;
- memilih teknologi berdasarkan preferensi tanpa mempertimbangkan kebutuhan;
- menganggap asumsi sebagai requirement final;
- mengubah banyak file tanpa scope task yang jelas;
- mengubah API, schema, atau data secara diam-diam;
- menyatakan test atau build berhasil tanpa menjalankannya;
- melupakan dokumentasi setelah implementasi;
- melakukan commit, push, deployment, atau tindakan destruktif tanpa izin;
- kehilangan konteks saat pekerjaan dilanjutkan pada sesi berikutnya.

Dengan protokol ini, proyek diarahkan melalui alur:

```text
Audit Repository
→ Business Discovery
→ Persetujuan BRD
→ Product Discovery
→ Persetujuan PRD
→ Technical Discovery
→ Architecture dan Perencanaan
→ Persetujuan Rencana
→ Implementasi per Task
→ Verifikasi
→ Release
→ Maintenance
```

---

## Konsep Utama

### 1. `AGENTS.md` adalah aturan permanen

`AGENTS.md` tidak hanya digunakan sekali saat setup. File ini tetap berada di root repository dan berlaku sepanjang umur proyek.

Yang biasanya dilakukan satu kali adalah bootstrap lengkap:

```text
Business Discovery
→ BRD
→ Product Discovery
→ PRD
→ Technical Discovery
→ Architecture
→ Implementation Plan
```

Setelah bootstrap selesai, sesi berikutnya dilanjutkan dari state dan tracking proyek tanpa mengulang wawancara penuh.

### 2. Kebutuhan mendahului teknologi

Urutan pengambilan keputusan adalah:

```text
Mengapa proyek dibangun?  → BRD
Apa yang harus dibangun?  → PRD
Bagaimana cara membangun? → Architecture dan Technical Plan
```

AI tidak boleh menetapkan teknologi final sebelum kebutuhan produk cukup jelas, kecuali melakukan feasibility investigation terbatas yang telah disetujui.

### 3. Implementasi dilakukan per task

Setelah rencana disetujui, AI tidak langsung mengerjakan seluruh aplikasi. Pekerjaan dibagi menjadi phase dan task kecil yang:

- memiliki requirement terkait;
- memiliki acceptance criteria;
- mempunyai dependency yang jelas;
- dapat diuji secara independen;
- meminta persetujuan sebelum dimulai;
- menghasilkan laporan dan bukti verifikasi.

### 4. Bukti lebih penting daripada klaim

Status `PASS` hanya boleh digunakan jika command atau pemeriksaan benar-benar dijalankan dan berhasil. Jika tidak dapat dijalankan, AI harus menggunakan status seperti `BLOCKED` atau `NOT RUN` beserta alasannya.

### 5. Dokumentasi adalah bagian dari pekerjaan

Setelah setiap task, AI harus memperbarui state, task list, report, traceability, dan dokumen lain yang terdampak. Dokumentasi bukan pekerjaan terpisah yang ditunda sampai akhir proyek.

### 6. Keamanan menyeluruh, implementasi berbasis risiko

Semua proyek wajib mengevaluasi seluruh kategori ancaman, tetapi kedalaman kontrol, dokumentasi, dan pengujian disesuaikan dengan risiko:

```text
Comprehensive Threat Coverage
+ Risk-Based Controls
+ Defense in Depth
+ Evidence-Based Verification
```

Keamanan mencakup `Identify/Govern`, `Protect`, `Detect`, `Respond`, dan `Recover`, bukan hanya secure coding. Kontrol yang tidak relevan harus diberi alasan `NOT APPLICABLE`; kontrol yang ditunda atau risikonya diterima membutuhkan owner, batas waktu/expiry, dan persetujuan.

---

## Kapan Digunakan

Gunakan `AGENTS.md` untuk:

- proyek baru yang baru berupa ide;
- repository baru yang hanya memiliki skeleton;
- existing repository yang belum terdokumentasi;
- proyek yang sudah berjalan dan perlu dilanjutkan;
- penambahan fitur;
- perbaikan bug;
- refactoring;
- perubahan database atau API;
- perubahan arsitektur;
- integrasi pihak ketiga;
- security hardening;
- persiapan release;
- deployment dan maintenance.

Kedalaman proses harus disesuaikan dengan kompleksitas proyek. Landing page sederhana tidak perlu diperlakukan seperti sistem finansial kritis, tetapi risiko penting tetap tidak boleh diabaikan.

---

## Persiapan

### Prasyarat

- Sebuah folder proyek atau repository.
- AI coding agent yang dapat membaca file di workspace.
- `AGENTS.md` ditempatkan di root proyek.
- Pengguna atau stakeholder yang dapat menjawab pertanyaan dan memberi persetujuan.

### Struktur awal minimum

```text
nama-proyek/
└── AGENTS.md
```

### Menyalin file ke proyek

Contoh konseptual:

```bash
mkdir nama-proyek
cp /lokasi/template/AGENTS.md nama-proyek/AGENTS.md
cd nama-proyek
```

Untuk existing repository:

```bash
cp /lokasi/template/AGENTS.md /lokasi/existing-repository/AGENTS.md
cd /lokasi/existing-repository
```

> Periksa terlebih dahulu apakah repository sudah memiliki `AGENTS.md`. Jangan menimpa instruksi lokal tanpa memahami dan menyelaraskan isinya.

### Apakah `AGENTS.md` dibaca otomatis?

Hal ini bergantung pada AI coding agent yang digunakan. Sebagian agent membaca `AGENTS.md` secara otomatis, tetapi sebagian lain mungkin memerlukan instruksi eksplisit.

Untuk memastikan file digunakan, awali sesi dengan:

```text
Baca dan ikuti AGENTS.md sebelum melakukan tindakan lain.
```

Jika platform memiliki hierarchy atau format instruksi sendiri, aturan platform yang lebih tinggi tetap berlaku.

---

## Quick Start

### Proyek baru

Gunakan prompt:

```text
Baca dan ikuti AGENTS.md.

Ini adalah proyek baru. Lakukan audit repository secara read-only,
tentukan kondisi proyek, lalu mulai Business Discovery secara bertahap.

Jangan membuat source code aplikasi sebelum BRD, PRD, arsitektur, dan
implementation plan disetujui sesuai gerbang yang ditentukan.
```

### Existing repository

Gunakan prompt:

```text
Baca dan ikuti AGENTS.md.

Repository ini sudah memiliki source code tetapi belum diinisialisasi
menggunakan standar proyek. Audit repository secara read-only.

Identifikasi stack, entry point, arsitektur, database, kontrak API,
test, CI/CD, deployment, dokumentasi, Git status, inkonsistensi, dan
informasi yang belum diketahui. Jangan mengubah file terlebih dahulu.
```

### Melanjutkan proyek

Gunakan prompt:

```text
Baca AGENTS.md dan lanjutkan proyek dari state saat ini.

Baca PROJECT-STATE, TASKS, laporan terbaru, decisions, open questions,
dokumentasi terkait, source code, serta Git status/diff.

Rangkum kondisi aktual, task aktif, task terakhir, blocker,
inkonsistensi, dan rekomendasi task berikutnya. Jangan mengubah atau
menjalankan kode sebelum task dan code plan disetujui.
```

---

## Penggunaan untuk Proyek Baru

### Langkah 1 — Audit awal

AI harus memeriksa isi repository secara read-only dan menentukan bahwa proyek masih baru.

Contoh respons yang diharapkan:

```markdown
Saya telah membaca AGENTS.md dan melakukan audit awal secara read-only.

**Kondisi proyek:** Proyek baru.
**Fakta:** Repository hanya berisi AGENTS.md dan belum memiliki source code.
**Status:** BUSINESS_DISCOVERY.
**Tahap berikutnya:** Menggali kebutuhan bisnis sebagai dasar BRD.
```

Pada tahap ini AI tidak boleh langsung membuat `package.json`, `go.mod`, `src/`, database, Dockerfile, atau source code aplikasi lainnya.

### Langkah 2 — Business Discovery

AI mengajukan 3–7 pertanyaan yang saling berkaitan per putaran, misalnya:

1. Masalah bisnis apa yang ingin diselesaikan?
2. Siapa yang mengalami masalah tersebut?
3. Bagaimana proses berjalan saat ini?
4. Apa dampak waktu, biaya, risiko, atau kualitasnya?
5. Hasil terukur apa yang ingin dicapai?

Pengguna menjawab berdasarkan kondisi aktual. Jika suatu jawaban belum diketahui, nyatakan dengan jujur agar dicatat sebagai pertanyaan terbuka, bukan ditebak menjadi requirement.

### Langkah 3 — Draft dan persetujuan BRD

Setelah informasi cukup, AI menyusun Business Requirements Document yang menjelaskan:

- latar belakang dan masalah;
- tujuan bisnis;
- stakeholder;
- proses As-Is dan To-Be;
- scope dan non-scope;
- business requirements;
- business rules;
- KPI;
- constraint, dependency, dan risiko;
- asumsi dan pertanyaan terbuka.

Pengguna kemudian memberi keputusan eksplisit:

```text
BRD disetujui sebagai dasar Product Discovery.
```

Atau:

```text
BRD belum disetujui. Perbaiki bagian scope dan KPI berikut: ...
```

### Langkah 4 — Product Discovery

Setelah BRD disetujui, AI menggali:

- product vision;
- persona dan role;
- access matrix;
- user journey dan use case;
- functional requirements;
- data requirements;
- UX/UI;
- integration;
- notification;
- analytics;
- nonfunctional requirements;
- error dan edge cases;
- prioritas dan release scope.

### Langkah 5 — Draft dan persetujuan PRD

PRD menjelaskan produk apa yang akan dibangun. Pengguna harus mengoreksi asumsi dan memberi persetujuan eksplisit:

```text
PRD disetujui sebagai dasar Technical Discovery.
```

### Langkah 6 — Technical Discovery

AI menggali kebutuhan teknis, termasuk:

- bahasa dan framework;
- runtime;
- database;
- API dan integrasi;
- authentication dan authorization;
- infrastruktur;
- deployment;
- observability;
- backup dan recovery;
- security dan test strategy;
- security risk tier (`LOW`, `MEDIUM`, `HIGH`, atau `CRITICAL`);
- asset, data flow, trust boundary, attack surface, dan threat coverage;
- preventive, detective, response, recovery control, telemetry, dan incident path;
- vulnerability management, residual risk, dan kebutuhan pengujian berizin;
- kemampuan tim dan biaya operasional.

Jika teknologi belum ditentukan, AI harus memberikan 2–4 opsi, trade-off, risiko, biaya, dan rekomendasi. Pengguna tetap menjadi pengambil keputusan.

Contoh jawaban:

```text
Saya memilih opsi kedua berdasarkan trade-off yang dijelaskan.
Catat keputusan ini dan lanjutkan penyusunan arsitektur.
```

### Langkah 7 — Dokumentasi teknis dan rencana

AI menyusun dokumentasi yang relevan, seperti:

- architecture;
- SOP spesifik proyek;
- security strategy;
- test strategy;
- implementation plan;
- phase detail;
- task list;
- traceability;
- project state.

Dokumen tambahan seperti API spec, ERD, threat model, atau deployment runbook harus diusulkan dan disetujui jika memang diperlukan.

### Langkah 8 — Persetujuan rencana

Pengguna meninjau stack, arsitektur, phase, task, test, security, deployment, dependency, dan risiko.

Contoh persetujuan:

```text
Arsitektur, SOP, implementation plan, phase, task list, test strategy,
security strategy, dan traceability disetujui.

Ubah proyek ke status READY dan ajukan task pertama beserta code plan.
Jangan mulai implementasi sebelum task dan code plan tersebut disetujui.
```

### Langkah 9 — Code plan dan implementasi per task

AI mengajukan detail task dan code plan melalui Gerbang 3. Plan menjelaskan kondisi implementasi saat ini, pendekatan langkah demi langkah, alur kode, file/symbol, kontrak/data/dependency, security, edge case, test, command, expected result, risiko, dan rollback. Setelah pengguna menyetujui keduanya, AI mengimplementasikan hanya scope dan plan tersebut, menjalankan verifikasi yang direncanakan, memperbarui tracking, membandingkan hasil aktual dengan expected result, lalu berhenti sebelum task berikutnya.

Jika ditemukan penyimpangan material, AI berhenti pada kondisi aman, menjelaskan dampak, merevisi plan, dan menunggu persetujuan ulang.

---

## Penggunaan untuk Existing Repository

Existing repository tidak boleh diperlakukan sebagai proyek kosong. Source code, schema, migration, konfigurasi, test, dan kontrak runtime aktual merupakan sumber kebenaran penting.

### Prosedur

1. Letakkan `AGENTS.md` di root repository.
2. Minta AI melakukan audit read-only.
3. AI mengidentifikasi stack dan struktur aktual.
4. AI membaca Git status/diff untuk mendeteksi pekerjaan yang belum selesai.
5. AI membandingkan kode, test, dan dokumentasi.
6. AI melaporkan inkonsistensi dan risiko.
7. AI hanya menanyakan informasi yang tidak dapat ditemukan.
8. BRD/PRD dan dokumentasi disusun berdasarkan kondisi aktual, bukan preferensi AI.
9. Perubahan atau eksekusi source code tetap menunggu persetujuan rencana proyek, task, dan code plan.

### Prompt audit existing repository

```text
Baca dan ikuti AGENTS.md.

Audit repository ini secara read-only. Jangan mengubah file.

Periksa struktur root, instruksi lokal, manifest dan lock file, source,
entry point, database dan migration, API contract, test, CI/CD, Docker
atau infrastructure as code, environment template, dokumentasi,
Git status, dan Git diff.

Klasifikasikan temuan sebagai fakta, inkonsistensi, risiko, constraint,
dependency, asumsi, dan pertanyaan terbuka. Tanyakan hanya informasi yang
tidak dapat ditemukan dari repository.
```

### Jangan lakukan ini

- Mengganti stack hanya karena AI lebih menyukai stack lain.
- Menulis ulang aplikasi sebelum memahami kontrak dan behavior aktual.
- Menghapus file yang tidak dikenal tanpa menyelidikinya.
- Menganggap dokumentasi lama selalu lebih benar daripada kode aktual.
- Menganggap kode aktual selalu benar jika bertentangan dengan requirement yang disetujui.

Jika kode dan dokumentasi berbeda, AI harus menunjukkan perbedaan, dampak, dan opsi penyelesaiannya sebelum mengubah kontrak, data, atau behavior.

---

## Melanjutkan Pekerjaan pada Sesi Berikutnya

`AGENTS.md` tetap digunakan pada setiap sesi, tetapi bootstrap penuh tidak diulang.

### File yang dibaca saat melanjutkan

```text
AGENTS.md
.agent/PROJECT-STATE.md
.agent/TASKS.md
.agent/REPORT.md
.agent/DECISIONS.md
.agent/OPEN-QUESTIONS.md
.agent/CHANGE-REQUESTS.md
docs/project/
source code aktual
Git status dan diff
```

### Tujuan pemeriksaan awal sesi

AI harus menentukan:

- status proyek;
- phase aktif;
- task aktif;
- task terakhir yang selesai dan terverifikasi;
- dependency yang sudah atau belum terpenuhi;
- perubahan yang belum tercatat;
- blocker;
- inkonsistensi dokumentasi dan kode;
- task berikutnya yang siap.

### Prompt sesi lanjutan

```text
Baca AGENTS.md dan lanjutkan proyek dari state saat ini.

Lakukan pemeriksaan awal sesi secara read-only. Baca PROJECT-STATE,
TASKS, laporan terbaru, decisions, open questions, change requests,
phase aktif, dokumentasi terkait, source code, serta Git status/diff.

Rangkum fakta, task aktif, task terakhir, blocker, inkonsistensi, risiko,
dan rekomendasi langkah berikutnya. Jangan mengubah atau menjalankan kode
sebelum task dan code plan disetujui.
```

### Jika task sebelumnya sudah selesai

AI memilih task berikutnya yang dependency-nya telah terpenuhi, menyusun code plan berdasarkan audit read-only, lalu mengajukan task dan plan tersebut melalui Gerbang 3.

### Jika task sebelumnya masih aktif

AI harus:

1. membaca scope dan acceptance criteria task;
2. memeriksa perubahan yang sudah dibuat;
3. memeriksa test atau verifikasi terakhir;
4. mengidentifikasi pekerjaan yang belum selesai;
5. memastikan scope tidak berubah;
6. melanjutkan task yang sama tanpa mengulang pekerjaan yang benar.

Persetujuan baru diperlukan jika ditemukan perubahan scope, risiko baru, migration tambahan, perubahan public contract, tindakan destruktif, atau pendekatan yang berbeda secara signifikan.

### Jika sesi terputus sebelum tracking diperbarui

AI harus merekonsiliasi:

```text
Git diff dan source code
↔ Scope task
↔ Acceptance criteria
↔ Test aktual
↔ Dokumentasi dan tracking
```

AI tidak boleh menganggap task selesai hanya karena source code sudah ada. Task baru dapat dinyatakan selesai setelah acceptance criteria dan verifikasi wajib terpenuhi.

---

## Tiga Gerbang Persetujuan

### Gerbang 1A — Persetujuan BRD

Memastikan masalah, tujuan, stakeholder, proses, scope, business rule, KPI, constraint, dependency, dan risiko telah benar.

Tanpa persetujuan ini, AI tidak boleh melanjutkan ke requirement produk final.

### Gerbang 1B — Persetujuan PRD

Memastikan pengguna, role, journey, fitur, data, UX, integrasi, NFR, edge case, prioritas, dan acceptance criteria telah benar.

Tanpa persetujuan ini, arsitektur dan implementation plan final tidak boleh ditetapkan.

### Gerbang 2 — Persetujuan teknis dan rencana

Mencakup:

- teknologi dan database;
- arsitektur;
- security, test, dan deployment strategy;
- security risk tier dan security owner;
- asset/data flow/trust boundary/attack surface;
- threat coverage dan trigger kontrol;
- prevention, detection, response, recovery, telemetry, dan incident path;
- security test, vulnerability management, deferred/accepted risk;
- phase dan dependency;
- task list;
- traceability;
- risiko dan trade-off.

Sebagai bagian Gerbang 2, **Security Planning Gate** harus memastikan seluruh kategori ancaman berstatus `APPLICABLE` atau `NOT APPLICABLE`; threat scenario/control/test/evidence yang relevan telah direncanakan; `DEFER`/`ACCEPT` hanya diterapkan pada risiko/kontrol spesifik; dan pengujian ofensif memiliki otorisasi tertulis serta Rules of Engagement.

Gerbang 2 untuk implementasi normal tidak boleh disetujui jika terdapat open `Critical` finding atau `Critical` residual risk, `High` finding/risk tidak memiliki owner dan rencana yang memenuhi policy, attack surface/data kritis belum diketahui, atau security verification wajib belum direncanakan.

Jika existing repository sudah memiliki kondisi `Critical`, task containment/remediation yang sempit dapat berjalan hanya melalui **Remediation Gate**. Persetujuan ini menetapkan `remediation_allowed = true` hanya untuk task yang disebutkan, sementara `implementation_allowed = false`; fitur normal dan release tetap terblokir sampai mitigasi diverifikasi dan kondisi `Critical` diturunkan.

Setelah Gerbang 2 disetujui, proyek dapat masuk status `READY` dan `implementation_allowed` menjadi `true`.

### Gerbang 3 — Persetujuan task dan code plan

Sebelum menulis, mengubah, menghasilkan, atau menjalankan kode, AI menjelaskan:

- tujuan, scope, dan out-of-scope;
- requirement, dependency, dan Definition of Ready;
- kondisi implementasi saat ini, gap, dan bukti terkait;
- pendekatan implementasi langkah demi langkah;
- alur input, validation/authorization, proses, output, side effect, error, dan recovery;
- file, component, serta function/class/module/symbol terdampak;
- API/interface/event/schema, database, migration, dependency, dan compatibility impact;
- data classification, trust boundary, attack surface, privilege, threat/control, telemetry, dan residual risk;
- edge case dan failure scenario;
- rencana test, expected result, dan evidence;
- command/tool, urutan, target/environment, klasifikasi `Read-only investigation`, `Non-mutating verification`, atau `Mutating/side-effecting`, serta efek samping;
- risiko, stop condition, rollback/roll-forward, dan cleanup.

Persetujuan hanya berlaku untuk task, file, symbol, behavior, dependency, command, environment, dan side effect yang dicantumkan. Jika terjadi penyimpangan material, AI harus berhenti, merevisi plan, dan meminta persetujuan ulang.

### Contoh persetujuan yang jelas

```text
Task dan code plan phase-02.3 disetujui sesuai scope, file/symbol,
command, verifikasi, risiko, dan rollback yang diajukan.

Kerjakan hanya plan tersebut. Jika perlu penyimpangan material, berhenti
pada kondisi aman dan minta persetujuan ulang. Setelah selesai, jalankan
verifikasi, bandingkan hasil dengan expected result, perbarui dokumentasi
dan tracking, lalu berhenti sebelum task berikutnya.
```

---

## Model Keamanan Menyeluruh

### Prinsip

Standar menggunakan model:

```text
Semua ancaman dievaluasi
→ Risiko diklasifikasikan
→ Kontrol relevan dipilih
→ Kontrol diimplementasikan
→ Kontrol diuji
→ Bukti dicatat
→ Risiko tersisa disetujui
→ Monitoring, response, dan recovery disiapkan
```

Tujuannya bukan menjanjikan bahwa peretasan mustahil terjadi, melainkan mengurangi peluang, membatasi dampak, mempercepat deteksi, menyediakan respons, dan membuktikan kontrol bekerja.

### Security risk tier

| Tingkat | Contoh karakteristik | Assurance minimum |
|---|---|---|
| `LOW` | Tanpa auth/data sensitif dan dampak rendah | Baseline, threat screening, secret/dependency scan sesuai stack |
| `MEDIUM` | Auth, database, data bisnis internal | Threat review, SAST/SCA, authorization test, vulnerability management |
| `HIGH` | PII, public API, multi-tenant, transaksi penting | Formal threat model, DAST/API test, infrastructure review, security gate |
| `CRITICAL` | Finansial, kesehatan, regulated, safety/high impact | Independent penetration test, adversarial exercise berizin, formal residual-risk review/approval, incident/recovery exercise |

Kompleksitas teknis dan risiko keamanan dinilai terpisah. Aplikasi kecil dapat tetap berisiko tinggi bila memproses credential atau data sensitif. Pada audit awal, tier bersifat provisional; tier final ditentukan setelah BRD, PRD, data classification, dan security discovery, lalu disetujui pada Security Planning Gate.

### Threat coverage matrix

Seluruh kategori berikut wajib dievaluasi:

| Kategori | Prevent | Detect | Respond/Recover | Verify |
|---|---|---|---|---|
| Reconnaissance | Kurangi exposure/metadata | Monitor probing | Tutup exposure | Attack-surface review |
| Scanning & Enumeration | Minimize service, rate limit | Scan/anomaly telemetry | Block/investigate | Authorized scanning |
| Vulnerability Assessment | Patch/hardening | Vulnerability monitoring | Remediation workflow | SAST/SCA/DAST/API/config scan |
| Exploitation | Secure design/coding, segmentation | Exploit telemetry | Isolate/remediate | Authorized security test |
| Password Attacks | Hashing, MFA, rate/step-up control | Login anomaly | Revoke/reset/contain | Controlled auth test |
| Social Engineering | Awareness dan verification SOP | Reporting channel | Account/vendor containment | Authorized simulation |
| Web/Application/API Attacks | OWASP-relevant controls | WAF/app/API telemetry | Block/patch | ASVS/WSTG/API testing |
| Privilege Escalation | Default deny, least privilege | Privileged-action alert | Revoke access | Authorization test |
| Lateral Movement | Segmentation/workload identity | East-west telemetry | Isolate workload | Authorized scenario |
| Persistence | Identity/key lifecycle, immutability | Integrity monitoring | Eradicate/rotate | Incident/detection exercise |
| Defense Evasion | Protected log/signed artifact | Security-control health alert | Restore telemetry | Detection validation |
| Exfiltration | Data/export/egress control | Bulk-access/egress alert | Revoke/contain/notify | Synthetic-data simulation |

Setiap kategori hanya menggunakan status:

- `APPLICABLE` — buat threat scenario `THR-*`, control, risk, test, dan evidence mapping;
- `NOT APPLICABLE` — alasan dan evidence konkret dicatat.

`DEFER` atau `ACCEPT` hanya boleh diterapkan pada `RISK-*`/`CTRL-*` spesifik, bukan seluruh kategori. Untuk setiap threat applicable, petakan preventive, detective, response/containment, dan recovery control beserta owner, test, serta evidence.

### Trigger kontrol

Kontrol tambahan menjadi wajib ketika proyek memiliki:

- authentication atau password;
- beberapa role/authorization;
- public web atau API;
- database atau data sensitif;
- PII, finansial, atau data kesehatan;
- file upload;
- cloud, container, atau IaC;
- multi-tenancy;
- CI/CD dan production deployment.

Contohnya, multi-tenant memicu tenant-isolation test; public API memicu object/function-level authorization dan resource-abuse test; local password authentication memicu hashing, reset/session security, serta brute-force/credential-stuffing defense. Authentication umum tetap memerlukan validasi identity provider/credential, token/session lifecycle, replay defense, anomaly monitoring, dan MFA sesuai risiko.

Selain 12 kategori utama, evaluasi availability/DoS, software/CI/CD supply chain, insider/authorized-user abuse, cryptographic/key compromise, destructive integrity/ransomware, dan third-party/provider compromise.

### Security assurance

Pemeriksaan dipilih berdasarkan stack dan risk tier:

```text
Repository/credential → Secret scanning
Dependency            → SCA/dependency audit
Source                 → SAST dan review
Web/API                → DAST/API security testing
Container              → Image scan
Infrastructure         → IaC/cloud/config scan
Access                 → Authorization/tenant-isolation test
Detection              → Expected-alert validation
High/Critical risk     → Independent penetration/adversarial exercise
```

Tool otomatis tidak membuktikan ketiadaan vulnerability. Penetration test juga tidak menggantikan secure design, review, monitoring, atau vulnerability management.

### Batas pengujian ofensif

Scanning aktif, exploitation, password simulation, phishing, privilege escalation, lateral movement, persistence, defense-evasion validation, dan exfiltration simulation selalu memerlukan otorisasi tertulis yang masih berlaku dari accountable asset/system owner. Kepemilikan tidak menggantikan otorisasi spesifik. Syaratnya mencakup:

- kewenangan approver dan terms/approval provider/customer/third party jika terdampak;
- target, account, environment, waktu, source, rate/volume, teknik boleh/dilarang yang jelas;
- Rules of Engagement;
- contact darurat, stop condition/kill switch;
- backup, rollback, dan cleanup;
- non-production representatif, test account, dan data sintetis/canary sebagai default;
- pengecualian production/real account tertulis dengan blast-radius limit dan real-time monitoring;
- evidence handling yang aman;
- legal/privacy/HR review jika relevan;
- expected alert dan success criteria.

Dilarang mengumpulkan password/OTP nyata, mengekstrak data production, menguji target di luar scope, melakukan destructive test, atau meninggalkan persistence tersembunyi. Jika otorisasi meragukan, pengujian harus berstatus `BLOCKED`.

### Vulnerability dan risk management

Bedakan project `security_risk_level`, `finding_severity`, dan `residual_risk_rating`; ketiganya bukan istilah yang sama. Setiap finding harus memiliki ID, affected asset, severity, evidence, exploitability, business impact, owner, target, remediation/mitigation, retest, dan closure evidence.

Open `Critical` finding dan `Critical` residual risk memblokir implementasi fitur normal dan release serta tidak dapat lolos hanya dengan risk acceptance. Hanya task containment/remediation yang sempit dan disetujui melalui Remediation Gate boleh berjalan. `High` finding/residual risk harus diperbaiki atau memiliki compensating control, owner, target, approver, dan expiry. Bedakan tool berhasil berjalan (`TOOL_RUN_PASS`) dari finding memenuhi policy gate (`POLICY_GATE_PASS`).

Security evidence wajib diklasifikasikan. Raw scanner output, exploit evidence, Rules of Engagement, dan vulnerability report sensitif tidak otomatis di-commit; simpan di storage terkontrol dan commit hanya ringkasan tersanitasi.

---

## Dokumentasi yang Dihasilkan

Struktur default bersifat adaptif:

```text
nama-proyek/
├── AGENTS.md
├── README.md
├── docs/
│   └── project/
│       ├── PROJECT-BRIEF.md
│       ├── BRD.md
│       ├── PRD.md
│       ├── ARCHITECTURE.md
│       ├── SOP.md
│       ├── IMPLEMENTATION-PLAN.md
│       ├── TEST-STRATEGY.md
│       ├── SECURITY.md
│       ├── TRACEABILITY.md
│       └── phases/
│           ├── phase-01-<nama>.md
│           ├── phase-02-<nama>.md
│           └── ...
└── .agent/
    ├── PROJECT-STATE.md
    ├── TASKS.md
    ├── REPORT.md
    ├── DECISIONS.md
    ├── OPEN-QUESTIONS.md
    └── CHANGE-REQUESTS.md
```

Tidak semua proyek memerlukan seluruh dokumen. AI harus menyesuaikan kedalaman dokumentasi dengan kompleksitas dan risiko.

Namun semua proyek tetap memiliki minimum security record. Untuk proyek sederhana, record dapat ditempatkan di `PROJECT-BRIEF.md` atau `SOP.md` dan memuat provisional/final risk tier, owner, threat-category disposition, applicable control/test, deferred/accepted risk spesifik, serta evidence reference. Proyek `MEDIUM` ke atas menggunakan `SECURITY.md` dan security traceability.

### Fungsi dokumen utama

| File | Fungsi |
|---|---|
| `AGENTS.md` | Aturan kerja permanen AI dan developer |
| `PROJECT-BRIEF.md` | Ringkasan singkat proyek |
| `BRD.md` | Alasan bisnis dan hasil yang ingin dicapai |
| `PRD.md` | Produk dan behavior yang harus dibangun |
| `ARCHITECTURE.md` | Cara teknis sistem dirancang |
| `SOP.md` | Prosedur kerja spesifik stack dan proyek |
| `IMPLEMENTATION-PLAN.md` | Phase, dependency, milestone, dan urutan implementasi |
| `TEST-STRATEGY.md` | Pendekatan verifikasi berbasis risiko |
| `SECURITY.md` | Kontrol keamanan dan privacy |
| `TRACEABILITY.md` | Hubungan requirement, keputusan, task, test, dan bukti |

### Fungsi file tracking `.agent/`

| File | Pertanyaan yang Dijawab |
|---|---|
| `PROJECT-STATE.md` | Proyek berada di status dan phase apa? |
| `TASKS.md` | Task mana yang siap, aktif, selesai, atau terblokir? |
| `REPORT.md` | Apa yang sudah dikerjakan dan dibuktikan? |
| `DECISIONS.md` | Keputusan penting apa yang diambil dan mengapa? |
| `OPEN-QUESTIONS.md` | Pertanyaan kritis apa yang belum terjawab? |
| `CHANGE-REQUESTS.md` | Perubahan scope apa yang diajukan dan bagaimana dampaknya? |

### Dokumen tambahan

Untuk proyek kompleks, AI dapat mengusulkan:

- API specification atau OpenAPI;
- ERD dan data dictionary;
- UX specification;
- threat model dan security test plan;
- attack-surface inventory;
- security control matrix;
- security risk register;
- Rules of Engagement untuk pengujian ofensif;
- compliance notes;
- data migration plan;
- deployment runbook;
- observability plan;
- backup/recovery plan;
- performance plan;
- release plan;
- operations manual;
- incident response;
- ADR.

Pembuatan dokumen tambahan harus memiliki alasan dan persetujuan karena setiap dokumen menambah biaya pemeliharaan.

---

## Siklus Pengerjaan Task

### Sebelum task

AI memeriksa Definition of Ready:

- tujuan jelas;
- requirement tersedia;
- acceptance criteria dapat diuji;
- scope dan out-of-scope jelas;
- dependency terpenuhi;
- kondisi implementasi saat ini, gap, dan bukti terkait dijelaskan;
- pendekatan implementasi langkah demi langkah tersedia;
- file, komponen, dan symbol terdampak diketahui;
- alur input, validation/authorization, proses, output, side effect, error, dan recovery dirancang;
- data, API/interface/event/schema, dependency, dan compatibility impact tersedia jika relevan;
- alternatif dependency/perubahan versi telah dievaluasi jika relevan;
- edge case dan failure scenario ditinjau;
- risiko utama diketahui;
- security/privacy impact telah ditinjau;
- perubahan data classification, trust boundary, attack surface, dan privilege ditinjau;
- threat, control, telemetry, security test, dan evidence dipetakan;
- residual risk memiliki owner/keputusan jika relevan;
- strategi test, expected result, dan evidence tersedia;
- command, target/environment, dan klasifikasi eksekusi dicantumkan;
- Rules of Engagement tersedia jika pengujian ofensif relevan;
- rollback/roll-forward/cleanup atau migration approach telah dipikirkan;
- tidak ada pertanyaan kritis terbuka;
- pengguna telah menyetujui task dan code plan.

### Tingkat code plan

| Aktivitas | Rencana Sebelum Eksekusi |
|---|---|
| Audit/investigasi read-only | Ringkas: tujuan, batas pemeriksaan, dan informasi yang dicari |
| Non-mutating verification: lint/typecheck/test/build/scan pasif | Command, tujuan, target/environment, expected result, dan konfirmasi tidak ada mutation yang diharapkan |
| Perubahan code/test/config/schema/dependency | Code plan lengkap dan persetujuan eksplisit |
| Mutating/side-effecting: migration/seed/install/generator/formatter-write/deployment | Code plan lengkap, side effect, blast radius, stop condition, backup, rollback/cleanup, dan persetujuan khusus jika berlaku |

Persetujuan code plan tidak menggantikan approval khusus untuk production, deployment, tindakan Git terbatas, destructive operation, atau offensive security testing.

### Saat task

AI harus:

- bekerja hanya dalam scope dan code plan yang disetujui;
- mengikuti pola dan arsitektur repository;
- mempertahankan backward compatibility kecuali perubahan disetujui;
- menambahkan atau memperbarui test;
- menangani error secara eksplisit;
- menjaga secret dan data sensitif;
- mencatat technical debt atau keputusan baru;
- berhenti dan meminta keputusan jika menemukan risiko penting atau penyimpangan material di luar plan.

Penyimpangan material mencakup file/symbol/environment baru, contract/schema/data/dependency baru, perubahan authorization/trust boundary, command atau side effect tambahan, strategi implementasi/test yang berubah, serta peningkatan risiko atau blast radius. Perubahan tersebut tidak boleh dijalankan sebelum plan direvisi dan disetujui ulang.

### Setelah task

AI harus:

1. membandingkan implementasi dengan code plan yang disetujui;
2. mencatat penyimpangan dan persetujuan ulang;
3. menjalankan verifikasi yang direncanakan;
4. mencatat command dan hasil aktual serta perbedaannya dari plan;
5. membandingkan hasil aktual dengan expected result;
6. memeriksa acceptance criteria dan rollback/cleanup plan;
7. memeriksa security, data, compatibility, dan regression impact;
8. memperbarui dokumentasi;
9. memperbarui traceability;
10. memperbarui task list;
11. menambahkan report;
12. memperbarui project state;
13. melaporkan risiko dan blocker;
14. merekomendasikan task berikutnya;
15. berhenti sampai ada persetujuan berikutnya.

Jika task mengubah attack surface atau risiko, AI juga memperbarui threat model, attack-surface inventory, control matrix, risk register, security test mapping, telemetry, dan incident/runbook yang terdampak.

### Status hasil verifikasi

| Status | Arti |
|---|---|
| `PASS` | Dijalankan dan berhasil |
| `FAIL` | Dijalankan dan gagal |
| `BLOCKED` | Tidak dapat dijalankan karena alasan konkret |
| `NOT RUN` | Belum dijalankan |
| `NOT APPLICABLE` | Tidak relevan untuk task tersebut |

### Kapan task dianggap selesai?

Task hanya boleh diberi status `✅ Selesai dan terverifikasi` jika implementasi sesuai code plan yang disetujui, penyimpangan material telah mendapat persetujuan ulang, hasil aktual telah dibandingkan dengan expected result, acceptance criteria terpenuhi, dan semua pemeriksaan wajib yang relevan telah lulus.

Jika implementasi sudah ditulis tetapi verifikasi wajib belum selesai, gunakan:

```text
🧪 Menunggu verifikasi
```

---

## Prompt Siap Pakai

### 1. Memulai proyek baru

```text
Baca dan ikuti AGENTS.md.

Inisialisasi proyek ini mulai dari audit repository secara read-only.
Setelah audit, lakukan Business Discovery secara bertahap untuk menyusun
BRD. Jangan membuat source code sebelum seluruh gerbang persetujuan yang
relevan selesai.
```

### 2. Mengaudit existing repository

```text
Baca dan ikuti AGENTS.md.

Audit repository secara read-only. Identifikasi stack, entry point,
arsitektur, database, migration, API, test, CI/CD, deployment,
dokumentasi, dan Git status/diff.

Laporkan fakta, inkonsistensi, informasi yang belum diketahui, risiko,
constraint, dan dependency. Jangan mengubah file.
```

### 3. Menyetujui BRD

```text
BRD disetujui sebagai dasar Product Discovery.
Catat persetujuan, perbarui status proyek, lalu mulai pertanyaan PRD
secara bertahap. Jangan masuk ke implementasi.
```

### 4. Menyetujui PRD

```text
PRD disetujui sebagai dasar Technical Discovery.
Mulai pembahasan teknologi, database, arsitektur, security, testing,
observability, dan deployment secara bertahap.
```

### 5. Meminta perbandingan teknologi

```text
Teknologi belum ditentukan. Berikan 2–4 opsi yang realistis berdasarkan
BRD dan PRD. Bandingkan kelebihan, kekurangan, risiko, biaya operasional,
kemampuan tim, dan kesesuaiannya. Berikan rekomendasi, tetapi jangan
menetapkan keputusan tanpa persetujuan saya.
```

### 6. Menyetujui rencana

```text
Arsitektur, teknologi, SOP, security strategy, test strategy,
implementation plan, phase, task list, dan traceability disetujui.

Ubah proyek ke status READY, lalu ajukan task pertama melalui Gerbang 3.
Jangan langsung mengimplementasikannya.
```

### 7. Meminta code plan sebelum eksekusi

```text
Untuk task phase-X.Y, susun code plan terperinci sebelum menulis,
mengubah, menghasilkan, atau menjalankan kode.

Jelaskan kondisi implementasi saat ini dan bukti terkait; tujuan, scope,
dan out-of-scope; pendekatan langkah demi langkah; alur input, validasi,
authorization, proses, output, side effect, error, dan recovery; file dan
symbol terdampak; API/data/schema/dependency/compatibility impact;
security/privacy/operational impact; edge case; test dan expected result;
command beserta klasifikasi Read-only investigation, Non-mutating
verification, atau Mutating/side-effecting; risiko, stop condition,
rollback/roll-forward, dan cleanup.

Jangan melakukan eksekusi sampai task dan code plan saya setujui.
```

### 8. Menyetujui task dan code plan

```text
Task dan code plan phase-X.Y disetujui sesuai scope, file/symbol,
behavior, command, verifikasi, risiko, dan rollback yang diajukan.

Kerjakan hanya plan tersebut. Jika perlu menyimpang secara material,
berhenti pada kondisi aman, jelaskan dampaknya, revisi plan, dan minta
persetujuan ulang. Setelah selesai, bandingkan hasil aktual dengan
expected result, perbarui dokumentasi dan tracking, lalu berhenti sebelum
task berikutnya.
```

### 9. Melanjutkan sesi

```text
Baca AGENTS.md dan lanjutkan proyek dari state saat ini.

Baca PROJECT-STATE, TASKS, REPORT terbaru, DECISIONS, OPEN-QUESTIONS,
CHANGE-REQUESTS, phase aktif, dokumentasi terkait, source code, dan Git
status/diff.

Rangkum kondisi aktual dan ajukan langkah berikutnya. Jangan mengubah
atau menjalankan kode sebelum task dan code plan disetujui.
```

### 10. Melanjutkan task aktif

```text
Lanjutkan task aktif sesuai AGENTS.md.

Validasi scope, code plan yang disetujui, acceptance criteria, perubahan
yang sudah ada, command/test yang telah dijalankan, penyimpangan, dan
pekerjaan yang belum selesai. Jangan mengulang pekerjaan yang benar,
menjalankan command di luar plan, atau berpindah ke task lain.
```

### 11. Audit status tanpa perubahan

```text
Audit state proyek secara read-only berdasarkan AGENTS.md.

Bandingkan PROJECT-STATE, TASKS, REPORT terbaru, requirement, phase aktif,
source code, test, dan Git diff. Laporkan inkonsistensi, blocker, risiko,
dan rekomendasi. Jangan mengubah file.
```

### 12. Meminta perubahan scope

```text
Ajukan permintaan berikut sebagai change request: <jelaskan perubahan>.

Analisis dampaknya terhadap BRD, PRD, arsitektur, database, API,
security, compatibility, testing, deployment, phase, task, dan timeline.
Jangan implementasikan sebelum perubahan requirement dan rencana
disetujui.
```

### 13. Melaporkan bug

```text
Ada bug: <jelaskan gejala>.

Ikuti AGENTS.md. Audit dan reproduksi masalah terlebih dahulu.
Dokumentasikan expected behavior, actual behavior, langkah reproduksi,
environment, severity, impact, evidence, kemungkinan komponen, dan
rencana regression test. Ajukan task perbaikan sebelum mengubah kode.
```

### 14. Mengevaluasi release

```text
Evaluasi release gate sesuai AGENTS.md.

Periksa scope, requirement, acceptance criteria, test, security scan,
migration, backup dan restore, rollback, environment, monitoring,
runbook, known issues, dan release notes.

Bedakan PASS, FAIL, BLOCKED, NOT RUN, dan NOT APPLICABLE. Jangan deploy.
```

### 15. Melakukan security planning review

```text
Lakukan Security Planning Gate sesuai AGENTS.md.

Tentukan security risk tier dan owner. Petakan asset, data flow, trust
boundary, attack surface, identity, public exposure, preventive control,
telemetry, response, recovery, dan security test.

Evaluasi reconnaissance, scanning/enumeration, vulnerability assessment,
exploitation, password attacks, social engineering, web/API attacks,
privilege escalation, lateral movement, persistence, defense evasion,
exfiltration, dan business-logic abuse.

Beri status APPLICABLE atau NOT APPLICABLE untuk setiap kategori.
Terapkan DEFER atau ACCEPT hanya pada RISK-*/CTRL-* spesifik dengan
owner, compensating control, target/review, approver, dan expiry sesuai
policy. Jangan melakukan pengujian ofensif.
```

### 16. Menyiapkan pengujian keamanan berizin

```text
Susun proposal security testing tanpa mengeksekusinya.

Dokumentasikan bukti otorisasi yang diperlukan, target, account,
environment, waktu, source, teknik yang diizinkan/dilarang, rate/volume,
data sintetis, expected alert, success criteria, contact darurat, kill
switch, backup, rollback, evidence handling, dan cleanup sebagai Rules
of Engagement.

Jika kepemilikan atau otorisasi target tidak jelas, tandai BLOCKED.
```

### 17. Meminta commit

```text
Periksa perubahan task aktif, status test, scope file, dan kemungkinan
secret. Jika aman dan seluruh perubahan sesuai scope, buat commit dengan
pesan yang mengikuti konvensi proyek. Jangan push.
```

Commit dan push adalah tindakan terpisah. Persetujuan commit tidak otomatis berarti persetujuan push.

---

## Bug, Perubahan Scope, dan Incident

### Bug

Bug harus memiliki:

- expected behavior;
- actual behavior;
- reproduction steps;
- environment dan version;
- severity dan impact;
- evidence;
- regression status;
- suspected component;
- verification plan.

AI sebaiknya menambahkan regression test yang gagal sebelum perbaikan dan lulus setelah perbaikan jika pendekatan tersebut sesuai dengan stack dan jenis bug.

### Perubahan scope

Fitur baru yang tidak terdapat dalam requirement disetujui harus diproses sebagai change request:

```text
Permintaan
→ Impact Analysis
→ Revisi BRD/PRD jika perlu
→ Revisi Architecture/Plan jika perlu
→ Persetujuan
→ Task baru
→ Gerbang 3
→ Implementasi
```

Jangan menyisipkan fitur besar ke task yang sedang aktif.

### Re-inisialisasi

Re-inisialisasi penuh hanya dipertimbangkan jika:

- tujuan produk berubah total;
- stack utama diganti;
- arsitektur dirombak;
- fungsi repository berubah;
- scope baru tidak lagi cocok dengan BRD/PRD;
- dokumentasi sangat tertinggal.

Re-inisialisasi memerlukan persetujuan eksplisit dan tidak boleh menghapus decision history.

### Production incident

Prioritas incident adalah:

1. melindungi pengguna dan data;
2. menstabilkan layanan;
3. mengumpulkan dan melindungi bukti secara aman;
4. melakukan containment minimum yang terukur;
5. mengomunikasikan status dan mengaktifkan escalation path;
6. mencabut/merotasi credential, session, atau key yang terdampak;
7. memeriksa lateral movement, persistence, defense evasion, dan exfiltration;
8. melakukan eradication dan recovery berdasarkan bukti;
9. memvalidasi monitoring dan indikator berulang;
10. melakukan root cause analysis;
11. membuat permanent fix serta regression/detection test;
12. memperbarui threat model, control matrix, risk register, dan post-incident review.

Emergency dapat mempercepat ceremony atau time window approval, tetapi minimum code plan tetap harus menjelaskan scope containment, command, target/environment, risiko, stop condition, verification, dan rollback/cleanup. Approval Git, deployment, destructive action, Remediation Gate, dan offensive security testing tetap terpisah. Emergency tidak mengizinkan destructive remediation tanpa persetujuan pihak yang berwenang, dan bukti tidak boleh dihapus hanya untuk mempercepat pemulihan.

---

## Release dan Deployment

Build yang berhasil tidak otomatis berarti aplikasi siap dirilis.

Sebelum release, periksa:

- scope release telah disetujui;
- requirement dan acceptance criteria terpenuhi;
- test wajib lulus;
- Security Planning Gate dan security verification wajib lulus;
- secret/SCA/SAST/DAST/API/container/IaC/config scan yang diwajibkan lulus;
- tidak ada open `Critical` finding atau `Critical` residual risk;
- `High` finding/residual risk diperbaiki atau memiliki keputusan risk treatment yang memenuhi policy;
- threat model, attack surface, control matrix, dan risk register mutakhir sesuai tier;
- authorization, tenant isolation, dan abuse-case test lulus jika relevan;
- security telemetry, tamper protection, alert, escalation, dan owner siap;
- incident containment, credential rotation, dan recovery path siap;
- migration diuji;
- backup tersedia dan restore telah diuji jika relevan;
- rollback atau roll-forward plan tersedia;
- environment variables dan secret siap serta dikelola aman;
- deferred/accepted risk belum kedaluwarsa dan memiliki owner;
- runbook tersedia;
- known issues terdokumentasi;
- release notes tersedia;
- approval release telah diberikan.

### Pemisahan approval

Approval berikut berbeda dan tidak saling menggantikan:

1. approval implementation plan;
2. approval task;
3. approval release;
4. approval deployment;
5. approval tindakan Git seperti commit atau push.

### Prompt deployment

Deployment hanya dilakukan setelah release gate diterima:

```text
Release gate telah ditinjau dan disetujui.
Lakukan deployment ke <environment> sesuai deployment runbook.

Jalankan pre-deployment check, migration sesuai urutan, smoke test,
post-deployment monitoring, dan laporkan hasil aktual. Jika terjadi
kegagalan yang memenuhi rollback criteria, berhenti dan ikuti runbook.
```

Jangan menaruh credential langsung di prompt, source code, atau dokumen yang akan di-commit. Gunakan secret manager atau mekanisme environment yang disetujui.

---

## Tanggung Jawab Pengguna dan AI

### Tanggung jawab pengguna

Pengguna atau stakeholder yang berwenang perlu:

- menjelaskan masalah dan tujuan bisnis;
- mengoreksi asumsi AI;
- menentukan stakeholder dan pengambil keputusan;
- menyetujui atau menolak BRD;
- menyetujui atau menolak PRD;
- memilih teknologi setelah memahami trade-off;
- menyetujui arsitektur dan implementation plan;
- menyetujui setiap task;
- menyediakan akses melalui mekanisme aman;
- menentukan penerimaan risiko;
- menyetujui release dan deployment;
- memberi instruksi eksplisit untuk commit atau push jika diperlukan.

### Tanggung jawab AI

AI harus:

- mempelajari repository sebelum bertanya atau mengubahnya;
- membedakan fakta, keputusan, asumsi, dan rekomendasi;
- bertanya secara bertahap;
- tidak memilih keputusan penting secara diam-diam;
- menjaga scope task;
- mempertimbangkan security, data, compatibility, dan operasi;
- menjalankan verifikasi aktual;
- melaporkan kegagalan dan blocker dengan jujur;
- memperbarui dokumentasi dan tracking;
- tidak membocorkan secret atau reasoning internal;
- berhenti sebelum task berikutnya sampai ada persetujuan.

### Bentuk jawaban persetujuan

Gunakan jawaban yang jelas:

```text
Disetujui.
```

```text
Disetujui dengan perubahan berikut: ...
```

```text
Belum disetujui karena: ...
```

```text
Tunda dan catat sebagai open question.
```

```text
Masukkan ke post-MVP, bukan scope release saat ini.
```

Hindari jawaban ambigu seperti “silakan lanjut” jika ada beberapa proposal atau task sekaligus. Sebutkan artefak atau task yang disetujui.

---

## Status dan Tracking

### Status proyek

| Status | Makna Ringkas |
|---|---|
| `UNINITIALIZED` | Informasi proyek belum cukup |
| `REPOSITORY_AUDIT` | Repository sedang diaudit |
| `BUSINESS_DISCOVERY` | Kebutuhan bisnis sedang digali |
| `BRD_PENDING_APPROVAL` | BRD menunggu persetujuan |
| `PRODUCT_DISCOVERY` | Kebutuhan produk sedang digali |
| `PRD_PENDING_APPROVAL` | PRD menunggu persetujuan |
| `TECHNICAL_DISCOVERY` | Solusi teknis sedang ditentukan |
| `DOCUMENTATION_DRAFT` | Dokumentasi teknis sedang disusun |
| `PLAN_PENDING_APPROVAL` | Rencana menunggu persetujuan |
| `READY` | Siap memulai task |
| `IN_PROGRESS` | Ada task aktif |
| `BLOCKED` | Pekerjaan terhalang |
| `ON_HOLD` | Proyek sengaja ditunda |
| `RELEASE_PENDING` | Menunggu release gate |
| `RELEASED` | Versi telah dirilis |
| `MAINTENANCE` | Dalam pemeliharaan |
| `COMPLETED` | Scope telah ditutup secara formal |

### Status task

| Status | Makna |
|---|---|
| `⬜ Belum siap` | Definition of Ready belum terpenuhi |
| `📋 Siap` | Dependency dan requirement siap |
| `🔄 Sedang dikerjakan` | Task aktif |
| `⚠️ Terblokir` | Ada blocker konkret |
| `🧪 Menunggu verifikasi` | Implementasi ada, verifikasi belum lengkap |
| `✅ Selesai dan terverifikasi` | Acceptance criteria dan quality gate terpenuhi |
| `⏸ Ditunda` | Sengaja ditunda |
| `❌ Dibatalkan` | Tidak akan dikerjakan |

### Traceability

Requirement idealnya dapat ditelusuri melalui rantai:

```text
OBJ
→ BUS
→ BRULE
→ FR/NFR
→ DEC/ADR
→ PHASE/TASK
→ TEST
→ EVIDENCE
```

Security menggunakan rantai tambahan:

```text
THREAT
→ CONTROL
→ RISK/FINDING
→ SECURITY TASK/TEST
→ EVIDENCE
→ RESIDUAL RISK
```

ID yang dapat digunakan: `THR-001`, `CTRL-001`, `RISK-001`, `FIND-001`, dan `STEST-001`.

Tujuannya agar setiap implementasi memiliki alasan, setiap requirement/threat memiliki control dan task, setiap task memiliki bukti, serta residual risk terlihat dan dimiliki pihak yang berwenang.

---

## Praktik yang Direkomendasikan

1. Simpan `AGENTS.md` di root repository.
2. Mulai setiap sesi dengan meminta AI membaca state dan diff aktual.
3. Gunakan satu task kecil dan fokus dalam satu siklus implementasi.
4. Berikan approval secara eksplisit dengan ID task.
5. Pastikan acceptance criteria dapat diuji.
6. Simpan keputusan penting dalam decision log.
7. Catat pertanyaan yang belum terjawab, jangan menebaknya.
8. Perbarui dokumen kontrak dalam task yang sama dengan perubahan kode.
9. Gunakan test data sintetis, bukan data production.
10. Periksa secret sebelum commit.
11. Bedakan approval task, commit, push, release, dan deployment.
12. Audit ulang state setelah sesi terputus atau berpindah agent.
13. Sesuaikan dokumentasi dengan risiko; jangan membuat dokumen kosong hanya untuk memenuhi struktur.
14. Gunakan feasibility spike jika ketidakpastian teknis tinggi, dengan scope dan waktu yang jelas.
15. Pertahankan riwayat keputusan ketika requirement berubah.
16. Tentukan security risk tier terpisah dari kompleksitas proyek.
17. Evaluasi seluruh threat category; jangan melewati kontrol tanpa alasan.
18. Hubungkan security requirement ke control, test, evidence, dan owner.
19. Lindungi log/audit dan validasi expected alert, bukan hanya menghasilkan log.
20. Jangan menjalankan security simulation tanpa otorisasi dan Rules of Engagement.

---

## Hal yang Dilarang atau Perlu Dihindari

- Meminta AI langsung membuat seluruh aplikasi dari deskripsi singkat.
- Menganggap rekomendasi AI otomatis menjadi keputusan.
- Menyetujui beberapa task tanpa scope yang jelas.
- Menyatakan task selesai tanpa test atau bukti yang diwajibkan.
- Menaruh password, token, API key, private key, atau credential di repository.
- Menggunakan data production sebagai fixture test tanpa perlindungan dan persetujuan.
- Mengubah public API atau schema secara diam-diam.
- Melakukan migration destruktif tanpa backup dan rollback/roll-forward plan.
- Menghapus pekerjaan yang belum dipahami.
- Menonaktifkan test atau security check hanya agar pipeline lulus.
- Melakukan commit, push, rebase, reset, clean, atau deployment tanpa instruksi eksplisit.
- Melakukan scanning aktif, exploitation, password/phishing simulation, lateral movement, persistence, defense-evasion, atau exfiltration test tanpa otorisasi dan Rules of Engagement.
- Menggunakan credential nyata atau data production sebagai artefak simulasi keamanan.
- Menutup finding tanpa retest/bukti atau menonaktifkan security check agar pipeline lulus.
- Menerima risiko kritis secara diam-diam atau tanpa pihak berwenang.
- Mengulang Business Discovery setiap sesi ketika requirement tidak berubah.
- Membuat dokumentasi enterprise yang tidak relevan untuk proyek sederhana.

---

## FAQ

### Apakah `AGENTS.md` hanya digunakan saat awal proyek?

Tidak. `AGENTS.md` berlaku sepanjang proyek. Yang hanya dilakukan di awal adalah bootstrap BRD, PRD, arsitektur, dan rencana. Sesi selanjutnya dilanjutkan dari `.agent/PROJECT-STATE.md` dan file tracking lainnya.

### Apakah semua sesi harus mengulang BRD dan PRD?

Tidak. BRD dan PRD hanya ditinjau ulang jika kebutuhan bisnis atau produk berubah, ditemukan konflik penting, atau dilakukan re-inisialisasi yang disetujui.

### Apakah AI boleh langsung mengerjakan task berikutnya setelah satu task selesai?

Tidak secara default. AI harus melaporkan hasil, memperbarui tracking, merekomendasikan task berikutnya, menyusun code plan berdasarkan audit read-only, lalu menunggu persetujuan Gerbang 3.

### Apakah persetujuan rencana berarti semua task otomatis disetujui?

Tidak. Persetujuan rencana membuka status `READY`, tetapi setiap task dan code plan tetap memerlukan persetujuan tersendiri.

### Apakah setiap command memerlukan code plan lengkap?

Tidak. Audit murni read-only cukup didahului tujuan dan batas pemeriksaan. Verifikasi non-mutating seperti lint, typecheck, test, atau build tetap harus dicantumkan dengan command, target, tujuan, dan expected result. Perubahan source/config/schema/dependency serta command mutating atau berefek samping memerlukan code plan lengkap dan persetujuan eksplisit.

### Bagaimana jika implementasi harus menyimpang dari code plan?

Jika penyimpangan mengubah file/symbol/environment, behavior, contract, data, dependency, security posture, command, side effect, risiko, atau rollback secara material, AI harus berhenti pada kondisi aman, menjelaskan dampak, merevisi plan, dan meminta persetujuan ulang. Perubahan kecil non-material cukup dicatat dalam laporan akhir.

### Apakah perubahan dokumentasi juga membutuhkan persetujuan task?

Koreksi typo atau perubahan dokumentasi kecil yang tidak mengubah requirement, kontrak, keputusan, atau behavior dapat dilakukan tanpa Gerbang 3 jika pengguna memang memintanya. Perubahan substantif tetap memerlukan analisis dan persetujuan.

### Bagaimana jika teknologi belum ditentukan?

AI harus menawarkan 2–4 opsi realistis, menjelaskan trade-off, memberikan rekomendasi, lalu menunggu keputusan pengguna.

### Bagaimana jika repository sudah memiliki teknologi tertentu?

AI harus menghormati stack aktual dan constraint repository. Mengganti stack memerlukan alasan, impact analysis, dan persetujuan eksplisit.

### Bagaimana jika kode dan dokumentasi bertentangan?

AI harus menunjukkan perbedaannya, menjelaskan dampak, mengklasifikasikan kemungkinan penyebab, dan meminta keputusan jika menyangkut kontrak, data, security, atau scope. Jangan langsung memilih salah satunya secara diam-diam.

### Apakah AI boleh membuat branch atau commit?

Hanya dengan instruksi eksplisit. Membaca Git status, diff, log, branch, dan tag diperbolehkan; tindakan yang mengubah history atau remote tidak boleh dilakukan otomatis.

### Apakah approval commit juga berarti approval push?

Tidak. Commit dan push adalah tindakan terpisah dan masing-masing membutuhkan instruksi eksplisit.

### Apakah AI boleh deploy setelah build berhasil?

Tidak. Build sukses bukan release approval. Release gate dan deployment approval tetap harus dilakukan.

### Bagaimana jika test tidak dapat dijalankan?

AI harus menjelaskan alasan konkret dan menggunakan status `BLOCKED` atau `NOT RUN`. Task tidak boleh dinyatakan selesai dan terverifikasi jika pemeriksaan wajib belum dilakukan.

### Apakah semua proyek membutuhkan semua dokumen?

Tidak. Dokumentasi harus adaptif. Proyek sederhana dapat menggunakan dokumen ringkas, sedangkan proyek kritis memerlukan kontrol dan dokumentasi lebih dalam.

### Bagaimana jika sesi AI berhenti di tengah task?

Pada sesi berikutnya, minta AI membaca state, task, report, source code, dan Git diff. AI harus merekonsiliasi pekerjaan aktual sebelum melanjutkan atau menyatakan task selesai.

### Bolehkah satu approval mencakup beberapa task?

Secara standar, satu task membutuhkan satu approval agar scope dan bukti tetap jelas. Jika pengguna ingin batching, risiko dan batas setiap task harus tetap eksplisit; perubahan ini sebaiknya dicatat sebagai keputusan workflow proyek.

### Di mana secret harus disimpan?

Gunakan environment variable, secret manager, atau mekanisme platform yang disetujui. `.env.example` hanya berisi nama variabel dan placeholder aman, bukan nilai rahasia.

---

## Troubleshooting

### AI langsung membuat kode tanpa discovery

Hentikan perubahan dan gunakan prompt:

```text
Hentikan implementasi. Baca AGENTS.md dan kembali ke state yang benar.
Audit repository secara read-only, identifikasi file yang sudah berubah,
dan laporkan dampaknya. Jangan membuang perubahan tanpa persetujuan.
```

### AI mengulang wawancara dari awal

Gunakan:

```text
Jangan ulangi discovery penuh. Baca PROJECT-STATE, BRD, PRD, TASKS,
REPORT, DECISIONS, dan OPEN-QUESTIONS. Tanyakan hanya informasi yang
belum tersedia atau yang berubah.
```

### State mengatakan selesai, tetapi test belum ada

Gunakan:

```text
Audit inkonsistensi antara PROJECT-STATE, TASKS, REPORT, source code,
dan bukti test. Jangan menganggap task selesai sebelum quality gate yang
wajib benar-benar dijalankan.
```

### AI ingin mengganti stack existing repository

Gunakan:

```text
Pertahankan stack repository saat ini. Jika penggantian memang perlu,
buat proposal terpisah berisi alasan, opsi, biaya migrasi, compatibility,
risiko, rollback, dan dampak terhadap timeline. Jangan mengubah stack.
```

### Ada perubahan pengguna yang tidak dikenali AI

Gunakan:

```text
Perubahan tersebut adalah pekerjaan pengguna. Jangan hapus atau timpa.
Audit dan jelaskan interaksinya dengan task aktif, lalu usulkan cara
melanjutkan dengan perubahan minimum.
```

### Dokumentasi tertinggal dari kode

Minta audit read-only terlebih dahulu. Tentukan apakah dokumentasi usang, kode tidak sesuai requirement, atau keduanya memerlukan rekonsiliasi. Perubahan kontrak atau requirement harus tetap disetujui.

### Task terlalu besar

Minta AI memecahnya:

```text
Task ini terlalu besar. Pecah menjadi task kecil yang fokus, dapat diuji,
memiliki dependency dan acceptance criteria sendiri, serta dapat direview
secara independen. Jangan mulai implementasi.
```

### Browser atau preview tidak tersedia

Preview visual bukan bukti build atau test. AI masih dapat melakukan pemeriksaan struktur file dan command lokal. Untuk UI, catat bahwa visual verification belum dilakukan dan jalankan ketika environment tersedia.

---

## Quick Reference

### Alur proyek baru

```text
1. Letakkan AGENTS.md di root
2. Buka AI agent dari root proyek
3. Audit repository secara read-only
4. Business Discovery
5. Draft dan approval BRD
6. Product Discovery
7. Draft dan approval PRD
8. Technical Discovery
9. Pilih teknologi
10. Susun architecture, SOP, security, dan test strategy
11. Susun implementation plan, phase, task, dan traceability
12. Approval Gerbang 2
13. Status READY
14. Audit read-only untuk task
15. Ajukan task dan code plan melalui Gerbang 3
16. Approval task dan code plan
17. Implementasi satu task sesuai plan
18. Berhenti dan re-approval jika ada penyimpangan material
19. Verifikasi dan bandingkan dengan expected result
20. Perbarui dokumentasi dan tracking
21. Berhenti sebelum task berikutnya
```

### Alur sesi lanjutan

```text
1. Baca AGENTS.md
2. Baca PROJECT-STATE
3. Baca TASKS dan REPORT terbaru
4. Baca DECISIONS dan OPEN-QUESTIONS terkait
5. Periksa source code dan Git diff
6. Identifikasi task aktif atau task siap
7. Rekonsiliasi inkonsistensi
8. Susun atau validasi code plan
9. Minta approval jika task/code plan belum disetujui atau ada penyimpangan material; untuk task aktif dengan approval valid, verifikasi approval tersebut
10. Kerjakan hanya scope dan plan yang disetujui
11. Verifikasi, bandingkan expected/actual, dan perbarui tracking
12. Berhenti
```

### Alur perubahan scope

```text
Request
→ Change Request
→ Impact Analysis
→ Revisi Requirement/Plan
→ Approval
→ Task dan Code Plan
→ Gerbang 3
→ Approval
→ Implementasi sesuai Plan
```

### Alur release

```text
Scope Selesai
→ Verifikasi Requirement
→ Test dan Security Gate
→ Migration/Backup/Rollback Check
→ Runbook dan Monitoring Check
→ Release Approval
→ Deployment Approval
→ Deploy
→ Smoke Test
→ Post-Deploy Monitoring
```

### Prompt paling praktis untuk setiap sesi

```text
Baca AGENTS.md dan lanjutkan proyek dari state saat ini.
Audit state, task, laporan terakhir, keputusan, pertanyaan terbuka,
dokumentasi terkait, source code, serta Git status/diff.

Rangkum kondisi aktual, inkonsistensi, blocker, risiko, dan langkah
berikutnya. Jangan mengubah atau menjalankan kode sebelum task dan code
plan disetujui.
```

---

## Penutup

`AGENTS.md` membantu menjadikan kerja dengan AI lebih terarah, dapat diaudit, dan aman. File ini tidak menggantikan keputusan manusia, review teknis, pengujian nyata, atau otorisasi organisasi. Perannya adalah memastikan bahwa kebutuhan, keputusan, implementasi, dan bukti tetap terhubung sepanjang siklus hidup proyek.

Untuk aturan lengkap, baca [`AGENTS.md`](./AGENTS.md).
