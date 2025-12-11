import type { Story } from "@ladle/react";
import { useState } from "react";
import { Slider } from "./Slider";

export const Default: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <div className="w-64">
      <Slider defaultValue={[50]} max={100} step={1} />
    </div>
  </div>
);

export const WithValue: Story = () => {
  const [value, setValue] = useState([50]);

  return (
    <div className="p-8 bg-gray-800 flex items-center justify-center">
      <div className="w-64 space-y-4">
        <div className="flex justify-between text-fg-primary text-sm">
          <span>音量</span>
          <span>{value[0]}%</span>
        </div>
        <Slider
          value={value}
          onValueChange={setValue}
          max={100}
          step={1}
        />
      </div>
    </div>
  );
};

export const Range: Story = () => {
  const [range, setRange] = useState([20, 80]);

  return (
    <div className="p-8 bg-gray-800 flex items-center justify-center">
      <div className="w-64 space-y-4">
        <div className="flex justify-between text-fg-primary text-sm">
          <span>范围</span>
          <span>{range[0]} - {range[1]}</span>
        </div>
        <Slider
          value={range}
          onValueChange={setRange}
          max={100}
          step={1}
        />
      </div>
    </div>
  );
};

export const AllSettings: Story = () => {
  const [volume, setVolume] = useState([70]);
  const [music, setMusic] = useState([50]);
  const [speed, setSpeed] = useState([3]);

  return (
    <div className="p-8 bg-gray-800 flex items-center justify-center">
      <div className="w-72 p-6 bg-surface-elevated rounded-lg space-y-6">
        <h3 className="text-fg-primary font-bold text-lg">游戏设置</h3>
        <div className="space-y-4">
          <div className="space-y-2">
            <div className="flex justify-between text-fg-primary text-sm">
              <span>音效音量</span>
              <span>{volume[0]}%</span>
            </div>
            <Slider value={volume} onValueChange={setVolume} max={100} step={1} />
          </div>
          <div className="space-y-2">
            <div className="flex justify-between text-fg-primary text-sm">
              <span>背景音乐</span>
              <span>{music[0]}%</span>
            </div>
            <Slider value={music} onValueChange={setMusic} max={100} step={1} />
          </div>
          <div className="space-y-2">
            <div className="flex justify-between text-fg-primary text-sm">
              <span>出牌速度</span>
              <span>{speed[0]}x</span>
            </div>
            <Slider value={speed} onValueChange={setSpeed} min={1} max={5} step={1} />
          </div>
        </div>
      </div>
    </div>
  );
};
