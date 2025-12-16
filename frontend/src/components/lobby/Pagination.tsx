import React from 'react';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { Button } from '../ui';
import { cn } from '@/lib/utils';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  totalPages,
  onPageChange
}) => {
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

  if (totalPages <= 1) {
    return null;
  }

  const visiblePages = getVisiblePages();

  return (
    <div className={cn(
      "flex items-center justify-center gap-1.5",
      "p-2 rounded-xl",
      "bg-surface-base/60 backdrop-blur-sm",
      "border border-stroke/30",
    )}>
      {/* 上一页按钮 */}
      <Button
        intent="neutral"
        size="sm"
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className={cn(
          "min-w-0 px-3",
          "border-stroke/50",
          !currentPage && "opacity-50",
        )}
      >
        <ChevronLeft className="w-4 h-4" />
        上一页
      </Button>

      {/* 页码 */}
      <div className="flex items-center gap-1 mx-2">
        {visiblePages.map((page, index) => (
          <React.Fragment key={index}>
            {page === '...' ? (
              <span className="px-2 py-1 text-fg-secondary/60 text-sm">...</span>
            ) : (
              <button
                onClick={() => onPageChange(page as number)}
                className={cn(
                  "min-w-[2rem] h-8 px-2 rounded-lg",
                  "text-sm font-medium",
                  "transition-all duration-fast",
                  currentPage === page
                    ? [
                        "bg-gradient-to-r from-[hsl(158,55%,35%)] to-[hsl(158,55%,42%)]",
                        "text-white",
                        "shadow-[var(--elevation-2),var(--shadow-relief)]",
                      ]
                    : [
                        "text-fg-secondary hover:text-fg-primary",
                        "hover:bg-surface-elevated/80",
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
          "min-w-0 px-3",
          "border-stroke/50",
          currentPage === totalPages && "opacity-50",
        )}
      >
        下一页
        <ChevronRight className="w-4 h-4" />
      </Button>
    </div>
  );
};

export default Pagination;
