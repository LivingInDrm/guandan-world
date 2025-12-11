import type { Story } from "@ladle/react";
import { Badge } from "./Badge";

export const Default: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <Badge>队友</Badge>
  </div>
);

export const AllVariants: Story = () => (
  <div className="p-8 bg-gray-800 flex flex-wrap gap-4">
    <div className="text-center">
      <Badge variant="landlord">地主</Badge>
      <p className="mt-2 text-white text-xs">landlord</p>
    </div>
    <div className="text-center">
      <Badge variant="farmer">农民</Badge>
      <p className="mt-2 text-white text-xs">farmer</p>
    </div>
    <div className="text-center">
      <Badge variant="teammate">队友</Badge>
      <p className="mt-2 text-white text-xs">teammate</p>
    </div>
    <div className="text-center">
      <Badge variant="owner">房主</Badge>
      <p className="mt-2 text-white text-xs">owner</p>
    </div>
  </div>
);

export const AllSizes: Story = () => (
  <div className="p-8 bg-gray-800 flex items-end gap-4">
    <div className="text-center">
      <Badge size="sm">小</Badge>
      <p className="mt-2 text-white text-xs">sm</p>
    </div>
    <div className="text-center">
      <Badge size="md">中</Badge>
      <p className="mt-2 text-white text-xs">md</p>
    </div>
  </div>
);

export const VariantSizeMatrix: Story = () => (
  <div className="p-8 bg-gray-800">
    <table className="border-separate border-spacing-4">
      <thead>
        <tr className="text-white text-xs">
          <th></th>
          <th>sm</th>
          <th>md</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td className="text-white text-xs pr-4">landlord</td>
          <td><Badge variant="landlord" size="sm">地主</Badge></td>
          <td><Badge variant="landlord" size="md">地主</Badge></td>
        </tr>
        <tr>
          <td className="text-white text-xs pr-4">farmer</td>
          <td><Badge variant="farmer" size="sm">农民</Badge></td>
          <td><Badge variant="farmer" size="md">农民</Badge></td>
        </tr>
        <tr>
          <td className="text-white text-xs pr-4">teammate</td>
          <td><Badge variant="teammate" size="sm">队友</Badge></td>
          <td><Badge variant="teammate" size="md">队友</Badge></td>
        </tr>
        <tr>
          <td className="text-white text-xs pr-4">owner</td>
          <td><Badge variant="owner" size="sm">房主</Badge></td>
          <td><Badge variant="owner" size="md">房主</Badge></td>
        </tr>
      </tbody>
    </table>
  </div>
);

export const InContext: Story = () => (
  <div className="p-8 bg-gray-800">
    <div className="flex items-center gap-3 bg-gray-700 p-4 rounded-lg">
      <div className="w-12 h-12 bg-gray-500 rounded-lg"></div>
      <div>
        <div className="flex items-center gap-2">
          <span className="text-white font-medium">玩家名称</span>
          <Badge variant="owner" size="sm">房主</Badge>
          <Badge variant="teammate" size="sm">队友</Badge>
        </div>
        <p className="text-gray-400 text-sm">等级 5</p>
      </div>
    </div>
  </div>
);
