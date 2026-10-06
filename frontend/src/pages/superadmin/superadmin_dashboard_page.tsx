export default function SuperadminDashboardPage() {
  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">Dashboard Superadmin</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="bg-gray-800 p-4 rounded-xl border border-gray-700">
          <p className="text-sm text-gray-400">Pengunjung Hari Ini</p>
          {/* <p className="text-3xl font-bold text-white mt-1">{mockVisitors.today} Orang</p> */}
        </div>
        <div className="bg-gray-800 p-4 rounded-xl border border-gray-700">
          <p className="text-sm text-gray-400">Pengunjung Minggu Ini</p>
          {/* <p className="text-3xl font-bold text-white mt-1">{mockVisitors.week} Orang</p> */}
        </div>
      </div>

      <div className="bg-gray-800 p-4 rounded-xl border border-gray-700">
        <h3 className="text-lg font-bold mb-3">Log Aktivitas Sistem</h3>
        <table className="w-full text-left text-sm text-gray-300">
          <thead className="bg-gray-700 text-gray-200">
            <tr>
              <th className="p-2">ID</th>
              <th className="p-2">Pengguna</th>
              <th className="p-2">Aktivitas</th>
              <th className="p-2">Waktu</th>
            </tr>
          </thead>
          <tbody>
            {/* {mockLogs.map((l) => (
              <tr key={l.id} className="border-b border-gray-700">
                <td className="p-2">{l.id}</td>
                <td className="p-2 text-white font-medium">{l.user}</td>
                <td className="p-2">{l.action}</td>
                <td className="p-2 text-gray-400">{l.time}</td>
              </tr>
            ))} */}
          </tbody>
        </table>
      </div>
    </div>
  )
}