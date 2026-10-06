export default function AdminSchedulesPage() {
  const handleAction = () => alert('Mock: Belum sambung API')

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-bold">Kelola Schedule</h2>
        <button onClick={handleAction} className="bg-red-600 px-3 py-1.5 rounded text-sm font-medium hover:bg-red-500">
          + Tambah Schedule
        </button>
      </div>
      <div className="bg-gray-800 rounded-xl overflow-hidden border border-gray-700">
        <table className="w-full text-left text-sm text-gray-300">
          <thead className="bg-gray-700 text-gray-200">
            <tr>
              <th className="p-3">ID</th>
              <th className="p-3">Film</th>
              <th className="p-3">Studio</th>
              <th className="p-3">Tanggal</th>
              <th className="p-3">Jam</th>
              <th className="p-3">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {/* {mockSchedules.map((s) => (
              <tr key={s.id} className="border-b border-gray-700">
                <td className="p-3">{s.id}</td>
                <td className="p-3 font-semibold text-white">{s.film}</td>
                <td className="p-3">{s.room}</td>
                <td className="p-3">{s.date}</td>
                <td className="p-3">{s.time}</td>
                <td className="p-3 space-x-2">
                  <button onClick={handleAction} className="text-blue-400 hover:underline">Edit</button>
                  <button onClick={handleAction} className="text-red-400 hover:underline">Hapus</button>
                </td>
              </tr>
            ))} */}
          </tbody>
        </table>
      </div>
    </div>
  )
}