import React from 'react';
import { ChevronLeft, ChevronRight, MoreHorizontal } from 'lucide-react';
import { Button } from '../ui';
import { cn } from '@/lib/utils';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  variant?: 'default' | 'arrows';
}

const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  totalPages,
  onPageChange,
  variant = 'default'
}) => {
  if (totalPages <= 1) {
    return null;
  }

  const hasPrev = currentPage > 1;
  const hasNext = currentPage < totalPages;

  if (variant === 'arrows') {
    return (
      <>
        {/* 左侧翻页按钮 */}
        <button
          onClick={() => hasPrev && onPageChange(currentPage - 1)}
          disabled={!hasPrev}
          className={cn(
            "absolute left-0 top-1/2 -translate-y-1/2 -translate-x-1/2 z-10",
            "w-8 h-8 mobile-landscape:w-6 mobile-landscape:h-6 rounded-full",
            "flex items-center justify-center",
            "bg-white border border-stroke/30",
            "shadow-[0_2px_8px_rgba(0,0,0,0.08)]",
            "transition-all duration-200",
            hasPrev
              ? "hover:bg-[hsl(40,8%,97%)] hover:border-state-active/50 hover:shadow-[0_4px_12px_rgba(0,0,0,0.12)] cursor-pointer"
              : "opacity-40 cursor-not-allowed",
          )}
        >
          <ChevronLeft className={cn(
            "w-5 h-5 mobile-landscape:w-4 mobile-landscape:h-4",
            hasPrev ? "text-fg-primary" : "text-fg-secondary",
          )} />
        </button>

        {/* 右侧翻页按钮 */}
        <button
          onClick={() => hasNext && onPageChange(currentPage + 1)}
          disabled={!hasNext}
          className={cn(
            "absolute right-0 top-1/2 -translate-y-1/2 translate-x-1/2 z-10",
            "w-8 h-8 mobile-landscape:w-6 mobile-landscape:h-6 rounded-full",
            "flex items-center justify-center",
            "bg-white border border-stroke/30",
            "shadow-[0_2px_8px_rgba(0,0,0,0.08)]",
            "transition-all duration-200",
            hasNext
              ? "hover:bg-[hsl(40,8%,97%)] hover:border-state-active/50 hover:shadow-[0_4px_12px_rgba(0,0,0,0.12)] cursor-pointer"
              : "opacity-40 cursor-not-allowed",
          )}
        >
          <ChevronRight className={cn(
            "w-5 h-5 mobile-landscape:w-4 mobile-landscape:h-4",
            hasNext ? "text-fg-primary" : "text-fg-secondary",
          )} />
        </button>
      </>
    );
  }

  const getVisiblePages = () => {
    const delta = 2;
    const range = [];
    const rangeWithDots = [];

    for (
      let i = Math.max(2, currentPage - delta);
      i <= Math.min(totalPages - 1, currentPage + delta);
      i++
    ) {
      range.push(i);
    }

    if (currentPage - delta > 2) {
      rangeWithDots.push(1, '...');
    } else {
      rangeWithDots.push(1);
    }

    rangeWithDots.push(...range);

    if (currentPage + delta < totalPages - 1) {
      rangeWithDots.push('...', totalPages);
    } else if (totalPages > 1) {
      rangeWithDots.push(totalPages);
    }

    return rangeWithDots;
  };

  const visiblePages = getVisiblePages();

  return (
    <div className={cn(
      "flex items-center justify-center gap-2 mobile-landscape:gap-1",
      "p-2.5 mobile-landscape:p-1.5 rounded-xl mobile-landscape:rounded-lg",
      "bg-white/70 backdrop-blur-sm",
      "border border-stroke/20",
      "shadow-[0_2px_8px_rgba(0,0,0,0.04),inset_0_1px_0_rgba(255,255,255,0.8)]",
    )}>
      {/* 上一页按钮 */}
      <Button
        intent="neutral"
        size="sm"
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className={cn(
          "min-w-0 px-3 mobile-landscape:px-2",
          "h-9 mobile-landscape:h-7",
          "bg-white hover:bg-[hsl(40,8%,97%)]",
          "border border-stroke/30 hover:border-stroke/50",
          "text-fg-secondary hover:text-fg-primary",
          "shadow-[0_1px_2px_rgba(0,0,0,0.04)]",
          "disabled:opacity-40 disabled:cursor-not-allowed",
          "mobile-landscape:text-xs",
          "transition-all duration-200",
        )}
      >
        <ChevronLeft className="w-4 h-4 mobile-landscape:w-3 mobile-landscape:h-3" />
        <span className="mobile-landscape:hidden">上一页</span>
      </Button>

      {/* 页码 */}
      <div className="flex items-center gap-1 mobile-landscape:gap-0.5 mx-1">
        {visiblePages.map((page, index) => (
          <React.Fragment key={index}>
            {page === '...' ? (
              <span className={cn(
                "flex items-center justify-center",
                "w-8 h-8 mobile-landscape:w-6 mobile-landscape:h-6",
                "text-fg-secondary/50",
              )}>
                <MoreHorizontal className="w-4 h-4 mobile-landscape:w-3 mobile-landscape:h-3" />
              </span>
            ) : (
              <button
                onClick={() => onPageChange(page as number)}
                className={cn(
                  "min-w-[2.25rem] mobile-landscape:min-w-[1.75rem]",
                  "h-9 mobile-landscape:h-7",
                  "px-2.5 mobile-landscape:px-1.5",
                  "rounded-lg mobile-landscape:rounded-md",
                  "text-sm mobile-landscape:text-xs font-semibold",
                  "transition-all duration-200",
                  currentPage === page
                    ? [
                        "bg-gradient-to-r from-[hsl(158,55%,38%)] to-[hsl(158,55%,45%)]",
                        "text-white",
                        "shadow-[0_2px_8px_hsla(158,55%,42%,0.3),inset_0_1px_0_rgba(255,255,255,0.15)]",
                        "scale-105",
                      ]
                    : [
                        "text-fg-secondary hover:text-fg-primary",
                        "hover:bg-[hsl(40,8%,95%)]",
                        "hover:scale-105",
                      ],
                )}
              >
                {page}
              </button>
            )}
          </React.Fragment>
        ))}
      </div>

      {/* 下一页按钮 */}
      <Button
        intent="neutral"
        size="sm"
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        className={cn(
          "min-w-0 px-3 mobile-landscape:px-2",
          "h-9 mobile-landscape:h-7",
          "bg-white hover:bg-[hsl(40,8%,97%)]",
          "border border-stroke/30 hover:border-stroke/50",
          "text-fg-secondary hover:text-fg-primary",
          "shadow-[0_1px_2px_rgba(0,0,0,0.04)]",
          "disabled:opacity-40 disabled:cursor-not-allowed",
          "mobile-landscape:text-xs",
          "transition-all duration-200",
        )}
      >
        <span className="mobile-landscape:hidden">下一页</span>
        <ChevronRight className="w-4 h-4 mobile-landscape:w-3 mobile-landscape:h-3" />
      </Button>
    </div>
  );
};

export default Pagination;
