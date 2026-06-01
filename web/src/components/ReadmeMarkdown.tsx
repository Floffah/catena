import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";

import ExpandableContent from "@/components/ExpandableContent";
import { cn } from "@/lib/utils";

export default function ReadmeMarkdown({
    content,
    className,
}: {
    content: string;
    className?: string;
}) {
    return (
        <ExpandableContent
            collapseLabel="Collapse README"
            expandLabel="Show full README"
        >
            <div
                className={cn(
                    "prose prose-sm max-w-none dark:prose-invert",
                    className,
                )}
            >
                <Markdown remarkPlugins={[remarkGfm]}>{content}</Markdown>
            </div>
        </ExpandableContent>
    );
}
