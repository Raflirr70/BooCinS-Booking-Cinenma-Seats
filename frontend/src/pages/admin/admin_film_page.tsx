export default function AdminFilmsPage() {
  const handleAction = () => alert('Mock: Belum sambung API')

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-bold">Kelola Film</h2>
        <button onClick={handleAction} className="bg-red-600 px-3 py-1.5 rounded text-sm font-medium hover:bg-red-500">
          + Tambah Film
        </button>
      </div>
      <div className="bg-gray-800 rounded-xl overflow-hidden border border-gray-700">
        <table className="w-full text-left text-sm text-gray-300">
          <thead className="bg-gray-700 text-gray-200">
            <tr>
              <th className="p-3">ID</th>
              <th className="p-3">Judul</th>
              <th className="p-3">Genre</th>
              <th className="p-3">Durasi</th>
              <th className="p-3">Harga</th>
              <th className="p-3">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {/* {mockFilms.map((f) => (
              <tr key={f.id} className="border-b border-gray-700">
                <td className="p-3">{f.id}</td>
                <td className="p-3 font-semibold text-white">{f.title}</td>
                <td className="p-3">{f.genre}</td>
                <td className="p-3">{f.duration}m</td>
                <td className="p-3">Rp {f.price.toLocaleString('id-ID')}</td>
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