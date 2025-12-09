import type { Story } from "@ladle/react";
import { useState } from "react";
import { Slider } from "./Slider";
import { Label } from "./Input";

export const Default: Story = () => {
  const [value, setValue] = useState([1.5]);

  return (
    <div className="p-8 max-w-md bg-table-100 rounded-lg">
      <div className="flex items-center gap-3">
        <Label className="shrink-0">缩放</Label>
        <Slider
          value={value}
          min={1}
          max={3}
          step={0.1}
          onValueChange={setValue}
          className="flex-1"
        />
      </div>
      <p className="mt-4 text-sm text-table-400">当前值: {value[0].toFixed(1)}</p>
    </div>
  );
};

export const Volume: Story = () => {
  const [volume, setVolume] = useState([50]);

  return (
    <div className="p-8 max-w-md bg-table-100 rounded-lg">
      <div className="flex items-center gap-3">
        <Label className="shrink-0">音量</Label>
        <Slider
          value={volume}
          min={0}
          max={100}
          step={1}
          onValueChange={setVolume}
          className="flex-1"
        />
      </div>
      <p className="mt-4 text-sm text-table-400">音量: {volume[0]}%</p>
    </div>
  );
};

export const Disabled: Story = () => (
  <div className="p-8 max-w-md bg-table-100 rounded-lg">
    <div className="flex items-center gap-3">
      <Label className="shrink-0">禁用状态</Label>
      <Slider
        value={[50]}
        min={0}
        max={100}
        disabled
        className="flex-1"
      />
    </div>
  </div>
);

export const WithoutLabel: Story = () => {
  const [value, setValue] = useState([30]);

  return (
    <div className="p-8 max-w-md bg-table-100 rounded-lg">
      <Slider
        value={value}
        min={0}
        max={100}
        step={5}
        onValueChange={setValue}
      />
      <p className="mt-4 text-sm text-table-400">值: {value[0]}</p>
    </div>
  );
};

export const Range: Story = () => {
  const [range, setRange] = useState([20, 80]);

  return (
    <div className="p-8 max-w-md bg-table-100 rounded-lg">
      <div className="flex items-center gap-3">
        <Label className="shrink-0">范围</Label>
        <Slider
          value={range}
          min={0}
          max={100}
          step={1}
          onValueChange={setRange}
          className="flex-1"
        />
      </div>
      <p className="mt-4 text-sm text-table-400">范围: {range[0]} - {range[1]}</p>
    </div>
  );
};
