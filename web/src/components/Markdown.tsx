import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import rehypeHighlight from 'rehype-highlight'

interface Props {
  children: string
}

/** Renders trusted markdown course content with grayscale code highlighting. */
export default function Markdown({ children }: Props) {
  return (
    <div className="prose-lua">
      <ReactMarkdown remarkPlugins={[remarkGfm]} rehypePlugins={[rehypeHighlight]}>
        {children}
      </ReactMarkdown>
    </div>
  )
}
