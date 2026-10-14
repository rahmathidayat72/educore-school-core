/**
 * Utility helper untuk ekspor dan impor data akademik ke/dari format Excel-compatible CSV dengan UTF-8 BOM
 * (Dapat langsung dibuka di Microsoft Excel, Google Sheets, LibreOffice Calc tanpa distorsi karakter)
 */

// Format cell value agar aman terhadap koma, tanda petik ganda, dan baris baru
const escapeCell = (val) => {
  if (val === null || val === undefined) return '""'
  const str = String(val).replace(/"/g, '""')
  return `"${str}"`
}

export function exportToCSV(filename, headers, rows) {
  if (!rows || rows.length === 0) {
    alert('Tidak ada data untuk diekspor.')
    return
  }

  // Baris header
  const headerLine = headers.map(h => escapeCell(h.label)).join(',')

  // Baris data
  const dataLines = rows.map((row, index) => {
    return headers.map(h => {
      let val
      if (h.key === '__index') {
        val = index + 1
      } else if (typeof h.formatter === 'function') {
        val = h.formatter(row)
      } else {
        val = row[h.key]
      }
      return escapeCell(val)
    }).join(',')
  })

  // Gabungkan dengan CRLF dan tambahkan UTF-8 Byte Order Mark (BOM: \uFEFF)
  const csvContent = '\uFEFF' + [headerLine, ...dataLines].join('\r\n')

  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.setAttribute('href', url)
  link.setAttribute('download', `${filename}_${new Date().toISOString().substring(0, 10)}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

/**
 * Mengunduh template CSV siap pakai untuk import data seragam
 * @param {string} filename - Nama file template
 * @param {Array<{key: string, label: string, example?: string}>} schema - Skema kolom template
 * @param {Array<Object>} sampleRows - Contoh baris data untuk panduan pengguna
 */
export function downloadCSVTemplate(filename, schema, sampleRows = []) {
  // Baris 1: Header kolom teknis/label
  const headerLine = schema.map(s => escapeCell(s.label || s.key)).join(',')

  // Baris 2+: Baris contoh data (dummy sample)
  const sampleLines = sampleRows.map(row => {
    return schema.map(s => escapeCell(row[s.key] !== undefined ? row[s.key] : (s.example || ''))).join(',')
  })

  const csvContent = '\uFEFF' + [headerLine, ...sampleLines].join('\r\n')
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.setAttribute('href', url)
  link.setAttribute('download', `${filename}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

/**
 * Parser file CSV sederhana & andal yang menangani baris baru, quote, dan koma RFC 4180
 * @param {string} csvText - Teks mentah dari file CSV
 * @returns {Array<Array<string>>}
 */
export function parseCSV(csvText) {
  // Bersihkan UTF-8 BOM jika ada
  let cleanText = csvText
  if (cleanText.charCodeAt(0) === 0xFEFF) {
    cleanText = cleanText.slice(1)
  }

  const rows = []
  let currentRow = []
  let currentVal = ''
  let insideQuotes = false

  for (let i = 0; i < cleanText.length; i++) {
    const char = cleanText[i]
    const nextChar = cleanText[i + 1]

    if (char === '"') {
      if (insideQuotes && nextChar === '"') {
        currentVal += '"'
        i++ // Lewati quote kedua
      } else {
        insideQuotes = !insideQuotes
      }
    } else if (char === ',' && !insideQuotes) {
      currentRow.push(currentVal.trim())
      currentVal = ''
    } else if ((char === '\r' || char === '\n') && !insideQuotes) {
      if (char === '\r' && nextChar === '\n') {
        i++ // Lewati \n pada CRLF
      }
      currentRow.push(currentVal.trim())
      if (currentRow.some(val => val !== '')) {
        rows.push(currentRow)
      }
      currentRow = []
      currentVal = ''
    } else {
      currentVal += char
    }
  }

  if (currentVal !== '' || currentRow.length > 0) {
    currentRow.push(currentVal.trim())
    if (currentRow.some(val => val !== '')) {
      rows.push(currentRow)
    }
  }

  return rows
}
