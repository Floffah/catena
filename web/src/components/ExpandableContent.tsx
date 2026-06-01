"use client";

import { IconChevronDown, IconChevronUp } from "@tabler/icons-react";
import { PropsWithChildren, useEffect, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface ExpandableContentProps extends PropsWithChildren {
    collapsedHeight?: number;
    expandLabel?: string;
    collapseLabel?: string;
    className?: string;
}

export default function ExpandableContent({
    children,
    collapsedHeight = 560,
    expandLabel = "Show full content",
    collapseLabel = "Collapse content",
    className,
}: ExpandableContentProps) {
    const contentRef = useRef<HTMLDivElement>(null);
    const [expanded, setExpanded] = useState(false);
    const [canExpand, setCanExpand] = useState(false);

    useEffect(() => {
        const content = contentRef.current;

        if (!content) {
            return;
        }

        const updateCanExpand = () => {
            setCanExpand(content.scrollHeight > collapsedHeight + 1);
        };

        updateCanExpand();

        const observer = new ResizeObserver(updateCanExpand);
        observer.observe(content);

        return () => {
            observer.disconnect();
        };
    }, [collapsedHeight]);

    return (
        <div className={cn("relative", className)}>
            <div
                ref={contentRef}
                className={cn(
                    "overflow-hidden transition-[max-height] duration-300 ease-out interpolate-allow",
                    expanded && "max-h-none",
                )}
                style={
                    expanded
                        ? { maxHeight: "max-content" }
                        : { maxHeight: collapsedHeight }
                }
            >
                {children}
            </div>

            {canExpand && (
                <div
                    className={cn(
                        "flex justify-center pt-4",
                        !expanded &&
                            "absolute inset-x-0 bottom-0 bg-linear-to-t from-card via-card/95 to-transparent pt-20",
                    )}
                >
                    <Button
                        type="button"
                        variant="secondary"
                        onClick={() => setExpanded((value) => !value)}
                    >
                        {expanded && (
                            <>
                                <IconChevronUp />
                                {collapseLabel}
                            </>
                        )}

                        {!expanded && (
                            <>
                                <IconChevronDown />
                                {expandLabel}
                            </>
                        )}
                    </Button>
                </div>
            )}
        </div>
    );
}
