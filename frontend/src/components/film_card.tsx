export interface ScheduleHome {
  schedule_id: number
  room_name: string
  date: string
  time: string
  available_seats: number
}

export interface FilmHome {
  film_id: number
  name: string
  cover: string
  duration: number
  price: number
  genres: string[]
  schedule: ScheduleHome[]
}

export default function FilmCard({ film }: { film: FilmHome }) {
  return (
    <div className="rounded-xl border border-gray-800 bg-gray-800/60 overflow-hidden flex flex-col">
      <img
        src={film.cover || '/assets/hero.png'}
        alt={film.name}
        className="h-64 w-full object-cover"
      />
      <div className="p-4 flex-1 flex flex-col justify-between">
        <div>
          <div className="flex items-center gap-2 mb-2 flex-wrap">
            {film.genres.map((g, idx) => (
              <span key={idx} className="text-xs bg-gray-700 text-gray-300 px-2 py-0.5 rounded">
                {g}
              </span>
            ))}
          </div>
          <h3 className="text-lg font-bold text-white mb-1">{film.name}</h3>
          <p className="text-xs text-gray-400">{film.duration} Menit</p>
        </div>
        <div className="mt-4 border-t border-gray-700/50 pt-3">
          <div className="flex justify-between items-center mb-2">
            <span className="text-sm font-bold text-red-400">
              Rp {film.price.toLocaleString('id-ID')}
            </span>
          </div>
          <p className="text-xs text-gray-400 mb-1">Jadwal Tayang:</p>
          {film.schedule.length > 0 ? (
            <div className="flex gap-1.5 flex-wrap">
              {film.schedule.map((s) => (
                <span
                  key={s.schedule_id}
                  className="text-xs bg-red-950 border border-red-800 text-red-200 px-2 py-1 rounded"
                >
                  {s.time} ({s.room_name})
                </span>
              ))}
            </div>
          ) : (
            <p className="text-xs text-gray-500 italic">Belum ada jadwal</p>
          )}
        </div>
      </div>
    </div>
  )
}