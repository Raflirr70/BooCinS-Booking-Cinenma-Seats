interface Promo {
  id: number
  title: string
  description: string
  img: string
}

export default function PromoBanner({ promo }: { promo: Promo | null }) {
  if (!promo) {
    return (
      <div className="rounded-xl border border-gray-800 bg-gray-800/50 p-6 text-center text-gray-400">
        Belum ada promo aktif saat ini.
      </div>
    )
  }

  return (
    <div className="rounded-2xl border border-red-900/50 bg-gradient-to-r from-red-950 to-gray-900 p-6 text-white shadow-xl flex flex-col md:flex-row items-center gap-6">
      {promo.img && (
        <img src={promo.img} alt={promo.title} className="w-full md:w-48 h-32 object-cover rounded-lg" />
      )}
      <div>
        <span className="text-xs uppercase tracking-widest text-red-400 font-bold">Promo Spesial</span>
        <h2 className="text-2xl font-bold mt-1">{promo.title}</h2>
        <p className="text-gray-300 text-sm mt-2">{promo.description}</p>
      </div>
    </div>
  )
}