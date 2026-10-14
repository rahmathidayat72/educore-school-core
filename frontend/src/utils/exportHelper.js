/**
 * Utility helper untuk ekspor data akademik ke format Excel-compatible CSV dengan UTF-8 BOM
 * (Dapat langsung dibuka di Microsoft Excel, Google Sheets, LibreOffice Calc tanpa distorsi karakter)
 */

export function exportToCSV(filename, headers, rows) {
  if (!rows || rows.length === 0) {
    alert('Tidak ada data untuk diekspor.')
    return
  }

  // Format cell value agar aman terhadap koma, tanda petik, dan baris baru
  const escapeCell = (val) => {
    if (val === null || val === undefined) return '""'
    const str = String(val).replace(/"/g, '""')
    return `"${str}"`
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
