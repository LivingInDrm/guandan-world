import type { Story } from "@ladle/react";
import { useState } from "react";
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogFooter,
  DialogTitle,
  DialogDescription,
} from "./Dialog";
import { Button } from "./Button";
import { Input, Label } from "./Input";

export const Default: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center min-h-[400px]">
    <Dialog>
      <DialogTrigger asChild>
        <Button>打开弹窗</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>游戏设置</DialogTitle>
          <DialogDescription>
            调整游戏规则和玩法参数。
          </DialogDescription>
        </DialogHeader>
        <div className="py-4">
          <p className="text-fg-primary">这是弹窗的内容区域。</p>
        </div>
      </DialogContent>
    </Dialog>
  </div>
);

export const WithFooter: Story = () => {
  const [open, setOpen] = useState(false);

  return (
    <div className="p-8 bg-gray-800 flex items-center justify-center min-h-[400px]">
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
          <Button>编辑资料</Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>编辑个人资料</DialogTitle>
            <DialogDescription>
              修改你的个人信息，点击保存完成更改。
            </DialogDescription>
          </DialogHeader>
          <div className="py-4 space-y-4">
            <div className="space-y-2">
              <Label htmlFor="nickname">昵称</Label>
              <Input id="nickname" placeholder="请输入昵称" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="signature">个人签名</Label>
              <Input id="signature" placeholder="请输入个人签名" />
            </div>
          </div>
          <DialogFooter>
            <Button intent="neutral" size="sm" onClick={() => setOpen(false)}>
              取消
            </Button>
            <Button intent="primary" size="sm" onClick={() => setOpen(false)}>
              保存
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export const LongContent: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center min-h-[400px]">
    <Dialog>
      <DialogTrigger asChild>
        <Button>查看规则</Button>
      </DialogTrigger>
      <DialogContent className="max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>游戏规则</DialogTitle>
          <DialogDescription>
            掼蛋游戏的完整规则说明。
          </DialogDescription>
        </DialogHeader>
        <div className="py-4 space-y-4">
          <h3 className="text-fg-primary font-semibold">一、基本规则</h3>
          <p className="text-fg-secondary text-sm">
            掼蛋是一种四人扑克牌游戏，使用两副牌共108张。四名玩家分为两组对抗，
            目标是尽快出完手中的牌。游戏以升级形式进行，从2开始，谁先打到A为胜。
          </p>
          <h3 className="text-fg-primary font-semibold">二、牌型介绍</h3>
          <p className="text-fg-secondary text-sm">
            单张：一张牌。对子：两张相同点数的牌。三张：三张相同点数的牌。
            炸弹：四张或更多相同点数的牌。顺子：五张或更多连续点数的牌。
          </p>
          <h3 className="text-fg-primary font-semibold">三、特殊规则</h3>
          <p className="text-fg-secondary text-sm">
            同花顺：五张或更多同花色的连续牌，可以打败普通炸弹。
            天王炸：四张王牌，是最大的牌型。级牌：当前等级的牌可以当作万能牌使用。
          </p>
          <h3 className="text-fg-primary font-semibold">四、得分规则</h3>
          <p className="text-fg-secondary text-sm">
            头游得3分，二游得2分，三游得1分，末游不得分。
            双上：同组两名玩家获得头游和二游，升3级。
            单下：对方一名玩家末游，升2级。
          </p>
          <h3 className="text-fg-primary font-semibold">五、游戏流程</h3>
          <p className="text-fg-secondary text-sm">
            1. 发牌：每人27张牌。2. 出牌：从有红桃2的玩家开始。
            3. 接牌：必须出比上家更大的牌型。4. 结束：第一个出完牌的玩家为头游。
          </p>
        </div>
      </DialogContent>
    </Dialog>
  </div>
);
