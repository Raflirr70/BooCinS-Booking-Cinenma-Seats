import { useState } from 'react'

export default function StaffCashierPage() {
//   const [selectedFilmId, setSelectedFilmId] = useState<number>(mockStaffFilms[0].id)
  const [qty, setQty] = useState<number>(1)
  const [cash, setCash] = useState<number>(0)

//   const film = mockStaffFilms.find((f) => f.id === selectedFilmId) || mockStaffFilms[0]
//   const totalPrice = film.price * qty
//   const change = cash - totalPrice

//   const handlePay = () => {
//     if (cash < totalPrice) {
//       alert('Uang cash kurang!')
//       return
//     }
//     alert(`Pembayaran Cash Berhasil! Kembalian: Rp ${change.toLocaleString('id-ID')}`)
//   }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <h2 className="text-2xl font-bold border-b border-gray-800 pb-3">Kasir Pembelian Tiket Offline</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-gray-800 p-5 rounded-xl border border-gray-700 space-y-4">
          <div>
            <label className="block text-sm text-gray-400 mb-1">Pilih Film</label>
            {/* <select value={selectedFilmId} onChange={(e) => setSelectedFilmId(Number(e.target.value))} className="w-full bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-white">
              {mockStaffFilms.map((f) => (
                <option key={f.id} value={f.id}>
                  {f.title} - Rp {f.price.toLocaleString('id-ID')}
                </option>
              ))}
            </select> */}
          </div>
          <div>
            <label className="block text-sm text-gray-400 mb-1">Jumlah Tiket</label>
            <input
              type="number"
              min={1}
              value={qty}
              onChange={(e) => setQty(Math.max(1, Number(e.target.value)))}
              className="w-full bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-white"
            />
          </div>
        </div>

        <div className="bg-gray-800 p-5 rounded-xl border border-gray-700 flex flex-col justify-between space-y-4">
          <div>
            <h3 className="text-sm font-semibold text-gray-400 uppercase">Ringkasan Pembayaran</h3>
            <p className="text-3xl font-bold text-red-500 mt-2">
              {/* Rp {totalPrice.toLocaleString('id-ID')} */}
            </p>
          </div>
          <div>
            <label className="block text-sm text-gray-400 mb-1">Uang Cash Diterima</label>
            <input
              type="number"
              value={cash}
              onChange={(e) => setCash(Number(e.target.value))}
              className="w-full bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-white"
              placeholder="Masukkan nominal"
            />
            {/* {cash > 0 && (
              <p className={`text-sm mt-2 ${change < 0 ? 'text-red-400' : 'text-green-400'}`}>
                {change < 0
                  ? `Kurang: Rp ${Math.abs(change).toLocaleString('id-ID')}`
                  : `Kembalian: Rp ${change.toLocaleString('id-ID')}`}
              </p>
            )} */}
          </div>
          <button
            // onClick={handlePay}
            className="w-full bg-red-600 hover:bg-red-500 py-3 rounded-lg font-bold text-white"
          >
            Bayar Cash
          </button>
        </div>
      </div>
    </div>
  )
}