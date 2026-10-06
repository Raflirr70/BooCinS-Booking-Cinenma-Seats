export default function AdminDashboardPage() {
  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">Dashboard Admin</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-gray-800 p-4 rounded-xl border border-gray-700">
          <p className="text-gray-400 text-sm">Total Film</p>
          {/* <p className="text-3xl font-bold mt-1">{mockFilms.length}</p> */}
        </div>
        <div className="bg-gray-800 p-4 rounded-xl border border-gray-700">
          <p className="text-gray-400 text-sm">Total Ruangan</p>
          {/* <p className="text-3xl font-bold mt-1">{mockRooms.length}</p> */}
        </div>
        <div className="bg-gray-800 p-4 rounded-xl border border-gray-700">
          <p className="text-gray-400 text-sm">Total Jadwal</p>
          {/* <p className="text-3xl font-bold mt-1">{mockSchedules.length}</p> */}
        </div>
        <div className="bg-gray-800 p-4 rounded-xl border border-gray-700">
          <p className="text-gray-400 text-sm">Promo Aktif</p>
          {/* <p className="text-3xl font-bold mt-1">{mockPromos.length}</p> */}
        </div>
      </div>
    </div>
  )
}
