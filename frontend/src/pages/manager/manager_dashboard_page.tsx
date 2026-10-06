export default function ManagerDashboardPage() {
  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">Dashboard Manager</h2>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-gray-800 p-4 rounded-xl border border-gray-700">
          <p className="text-sm text-gray-400">Pendapatan Online</p>
          <p className="text-2xl font-bold text-green-400 mt-1">
            {/* Rp {mockRevenue.online.toLocaleString('id-ID')} */}
          </p>
        </div>
        <div className="bg-gray-800 p-4 rounded-xl border border-gray-700">
          <p className="text-sm text-gray-400">Pendapatan Offline</p>
          <p className="text-2xl font-bold text-blue-400 mt-1">
            {/* Rp {mockRevenue.offline.toLocaleString('id-ID')} */}
          </p>
        </div>
        <div className="bg-gray-800 p-4 rounded-xl border border-gray-700">
          <p className="text-sm text-gray-400">Total Pendapatan</p>
          <p className="text-2xl font-bold text-red-500 mt-1">
            {/* Rp {mockRevenue.total.toLocaleString('id-ID')} */}
          </p>
        </div>
      </div>

      <div className="bg-gray-800 p-4 rounded-xl border border-gray-700">
        <h3 className="text-lg font-bold mb-3">Statistik Transaksi Booking</h3>
        <table className="w-full text-left text-sm text-gray-300">
          <thead className="bg-gray-700 text-gray-200">
            <tr>
              <th className="p-2">ID</th>
              <th className="p-2">Film</th>
              <th className="p-2">Jumlah Tiket</th>
              <th className="p-2">Tipe</th>
              <th className="p-2">Total Harga</th>
            </tr>
          </thead>
          <tbody>
            {/* {mockBookings.map((b) => (
              <tr key={b.id} className="border-b border-gray-700">
                <td className="p-2">{b.id}</td>
                <td className="p-2 text-white font-medium">{b.film}</td>
                <td className="p-2">{b.seats} Kursi</td>
                <td className="p-2">{b.type}</td>
                <td className="p-2">Rp {b.total.toLocaleString('id-ID')}</td>
              </tr>
            ))} */}
          </tbody>
        </table>
      </div>
    </div>
  )
}