import type { ReactNode } from "react";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "#/components/ui";
import { cn } from "#/lib/utils.ts";

export type TabItem = {
  value: string;
  label: ReactNode;
  /** Optional mono count badge next to the label. */
  count?: number | string;
  content: ReactNode;
  disabled?: boolean;
};

export type TabGroupProps = {
  value: string;
  onValueChange: (value: string) => void;
  items: TabItem[];
  className?: string;
  listClassName?: string;
  /** Applied to every tab panel. */
  contentClassName?: string;
};

// Underline tab strip + panels. Wraps the ui Tabs primitives with the line
// style used across app pages (members, vault detail, etc.).
export function TabGroup({
  value,
  onValueChange,
  items,
  className,
  listClassName,
  contentClassName,
}: TabGroupProps) {
  return (
    <Tabs
      value={value}
      onValueChange={onValueChange}
      className={cn("gap-5", className)}
    >
      <TabsList className={cn(listClassName)}>
        {items.map((item) => (
          <TabsTrigger key={item.value} value={item.value}>
            {item.label}
          </TabsTrigger>
        ))}
      </TabsList>

      {items.map((item) => (
        <TabsContent
          key={item.value}
          value={item.value}
          className={cn("flex flex-col gap-4", contentClassName)}
        >
          {item.content}
        </TabsContent>
      ))}
    </Tabs>
  );
}
