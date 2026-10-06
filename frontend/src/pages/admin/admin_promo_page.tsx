export default function AdminPromosPage() {
  const handleAction = () => alert('Mock: Belum sambung API')

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-bold">Kelola Promo</h2>
        <button onClick={handleAction} className="bg-red-600 px-3 py-1.5 rounded text-sm font-medium hover:bg-red-500">
          + Tambah Promo
        </button>
      </div>
      <div className="bg-gray-800 rounded-xl overflow-hidden border border-gray-700">
        <table className="w-full text-left text-sm text-gray-300">
          <thead className="bg-gray-700 text-gray-200">
            <tr>
              <th className="p-3">ID</th>
              <th className="p-3">Judul Promo</th>
              <th className="p-3">Deskripsi</th>
              <th className="p-3">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {/* {mockPromos.map((p) => (
              <tr key={p.id} className="border-b border-gray-700">
                <td className="p-3">{p.id}</td>
                <td className="p-3 font-semibold text-white">{p.title}</td>
                <td className="p-3">{p.desc}</td>
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