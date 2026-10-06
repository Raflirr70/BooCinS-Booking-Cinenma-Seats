import { useEffect, useState } from 'react'
import api, { apiMessage } from '../../lib/api'
import Navbar from '../../components/navbar'
import PromoBanner from '../../components/promo_banner'
import FilmCard, { type FilmHome } from '../../components/film_card'

interface HomeData {
  user_id: number
  name: string
  promo: { id: number; title: string; description: string; img: string } | null
  films: FilmHome[]
}

export default function HomePage() {
  const [data, setData] = useState<HomeData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    (async () => {
      try {
        const res = await api.get('/home')
        setData(res.data.data as HomeData)
      } catch (e) {
        setError(apiMessage(e))
      } finally {
        setLoading(false)
      }
    })()
  }, [])

  return (
    <div className="min-h-screen bg-gray-900 text-white flex flex-col">
      <Navbar />
      <main className="max-w-7xl mx-auto px-4 py-8 flex-1 w-full space-y-8">
        {loading && <p className="text-center text-gray-400 py-12">Loading data home...</p>}
        {error && <p className="text-center text-red-400 py-12">{error}</p>}
        {!loading && !error && data && (
          <>
            <section>
              <h2 className="text-xl font-bold mb-4 border-l-4 border-red-600 pl-3">Sedang Tayang</h2>
              {data.films.length === 0 ? (
                <p className="text-gray-400">Belum ada film yang tersedia saat ini.</p>
              ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
                  {data.films.map((film) => (
                    <FilmCard key={film.film_id} film={film} />
                  ))}
                </div>
              )}
            </section>
            <PromoBanner promo={data.promo} />
          </>
        )}
      </main>
    </div>
  )
}