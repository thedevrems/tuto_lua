import type { Lesson } from '../types'
import Markdown from './Markdown'

export default function LessonContent({ lesson }: { lesson: Lesson }) {
  return (
    <div className="h-full overflow-y-auto">
      <div className="max-w-3xl mx-auto px-8 py-10">
        <Markdown>{lesson.content}</Markdown>
      </div>
    </div>
  )
}
