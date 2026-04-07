import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import remarkFrontmatter from "remark-frontmatter";
import rehypeHighlight from "rehype-highlight";
import { Eye, Code, Copy, Download, Lock, Timer } from "lucide-react";
import { useToast } from "@/hooks/useToast";

interface NoteEditorProps {
  note: {
    id: string;
    title: string;
    content: string;
    note_type?: "normal" | "protected" | "self_destructing";
    self_destruct_at?: string | null;
  };
  onSave: (title: string, content: string) => Promise<void>;
  noteType?: "normal" | "protected" | "self_destructing";
}

export default function NoteEditor({
  note,
  onSave,
  noteType = "normal",
}: NoteEditorProps) {
  const [title, setTitle] = useState(note.title);
  const [content, setContent] = useState(note.content);
  // const [isSaving, setIsSaving] = useState(false);
  const [saveStatus, setSaveStatus] = useState<string | null>(null);
  const [viewMode, setViewMode] = useState<"split" | "editor" | "preview">(
    "split",
  );
  const { toast } = useToast();

  // Auto-save with debounce
  useEffect(() => {
    const timer = setTimeout(async () => {
      if (title !== note.title || content !== note.content) {
        // setIsSaving(true);
        setSaveStatus("Saving...");
        try {
          await onSave(title, content);
          setSaveStatus("Saved");
          setTimeout(() => setSaveStatus(null), 2000);
        } catch (error) {
          setSaveStatus("Error saving");
          console.error(error);
        }
        //  finally {
        //   // setIsSaving(false);
        // }
      }
    }, 1500);

    return () => clearTimeout(timer);
  }, [title, content, note, onSave]);

  const handleCopyMarkdown = async () => {
    try {
      await navigator.clipboard.writeText(content);
      toast({
        title: "Copied!",
        description: "Markdown content copied to clipboard",
      });
    } catch (error) {
      console.error(error);
      toast({
        title: "Error",
        description: "Failed to copy content",
        variant: "destructive",
      });
    }
  };

  const handleDownloadMarkdown = () => {
    const element = document.createElement("a");
    const file = new Blob([content], { type: "text/markdown" });
    element.href = URL.createObjectURL(file);
    element.download = `${title || "note"}.md`;
    document.body.appendChild(element);
    element.click();
    document.body.removeChild(element);
  };

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Header with title and controls */}
      <div className="border-b border-border p-4 flex flex-col gap-4">
        <div className="flex items-center justify-between gap-4">
          <div className="flex items-center gap-2 flex-1">
            <Input
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Note title..."
              className="text-lg font-semibold border-0 px-0 py-0 focus-visible:ring-0"
            />
            {noteType === "protected" && (
              <div className="flex items-center gap-1 px-2 py-1 bg-amber-100 dark:bg-amber-900/30 rounded text-xs text-amber-700 dark:text-amber-400 flex-shrink-0">
                <Lock className="w-3 h-3" />
                Protected
              </div>
            )}
            {noteType === "self_destructing" && (
              <div className="flex items-center gap-1 px-2 py-1 bg-red-100 dark:bg-red-900/30 rounded text-xs text-red-700 dark:text-red-400 flex-shrink-0">
                <Timer className="w-3 h-3" />
                Expires
              </div>
            )}
          </div>
          <div className="flex items-center gap-2">
            {saveStatus && (
              <span className="text-xs text-muted-foreground whitespace-nowrap">
                {saveStatus}
              </span>
            )}
          </div>
        </div>

        {/* View mode and action buttons */}
        <div className="flex items-center justify-between gap-2 flex-wrap">
          <div className="flex items-center gap-2">
            <Button
              variant={viewMode === "editor" ? "default" : "ghost"}
              size="sm"
              onClick={() => setViewMode("editor")}
              className="gap-2"
            >
              <Code className="w-4 h-4" />
              <span className="hidden sm:inline">Edit</span>
            </Button>
            <Button
              variant={viewMode === "preview" ? "default" : "ghost"}
              size="sm"
              onClick={() => setViewMode("preview")}
              className="gap-2"
            >
              <Eye className="w-4 h-4" />
              <span className="hidden sm:inline">Preview</span>
            </Button>
            <Button
              variant={viewMode === "split" ? "default" : "ghost"}
              size="sm"
              onClick={() => setViewMode("split")}
              className="gap-2"
            >
              <Code className="w-3 h-3" />
              <Eye className="w-3 h-3" />
              <span className="hidden sm:inline">Split</span>
            </Button>
          </div>

          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={handleCopyMarkdown}
              className="gap-2"
            >
              <Copy className="w-4 h-4" />
              <span className="hidden sm:inline">Copy</span>
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={handleDownloadMarkdown}
              className="gap-2"
            >
              <Download className="w-4 h-4" />
              <span className="hidden sm:inline">Download</span>
            </Button>
          </div>
        </div>
      </div>

      {/* Editor and Preview - DEFAULT SPLIT VIEW WITH REALTIME PREVIEW */}
      <div className="flex-1 overflow-hidden">
        {viewMode === "split" ? (
          <div className="grid grid-cols-2 h-full gap-0 divide-x divide-border">
            {/* Editor side */}
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="Write your markdown here...

# Heading 1
## Heading 2

**Bold** and *italic*

- Bullet points
- Another point

1. Numbered list
2. Second item

\`\`\`javascript
// Code block
const hello = 'world';
\`\`\`

> Blockquote

[Link](https://example.com)

| Table | Header |
|-------|--------|
| Cell  | Cell   |"
              className="w-full h-full p-6 bg-background text-foreground resize-none focus:outline-none font-mono text-sm"
            />
            {/* Preview side - REALTIME */}
            <div className="overflow-auto h-full p-6 bg-muted/30">
              <div className="prose prose-sm dark:prose-invert max-w-none">
                <ReactMarkdown
                  remarkPlugins={[remarkGfm, remarkFrontmatter]}
                  rehypePlugins={[rehypeHighlight]}
                >
                  {content || "Start typing to see preview..."}
                </ReactMarkdown>
              </div>
            </div>
          </div>
        ) : viewMode === "preview" ? (
          <div className="overflow-auto h-full p-6 bg-muted/30">
            <div className="prose prose-sm dark:prose-invert max-w-none">
              <ReactMarkdown
                remarkPlugins={[remarkGfm, remarkFrontmatter]}
                rehypePlugins={[rehypeHighlight]}
              >
                {content || "Start typing to see preview..."}
              </ReactMarkdown>
            </div>
          </div>
        ) : (
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="Write your markdown here...

# Heading 1
## Heading 2

**Bold** and *italic*

- Bullet points
- Another point

1. Numbered list
2. Second item

\`\`\`javascript
// Code block
const hello = 'world';
\`\`\`

> Blockquote

[Link](https://example.com)

| Table | Header |
|-------|--------|
| Cell  | Cell   |"
            className="w-full h-full p-6 bg-background text-foreground resize-none focus:outline-none font-mono text-sm"
          />
        )}
      </div>
    </div>
  );
}
