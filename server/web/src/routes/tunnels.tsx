import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table";
import { useToast } from "~/components/ui/use-toast";
import { request, Tunnel } from "~/lib/request";
import { useEffect, useState } from "react";
import {
  ArrowPathIcon,
  ClipboardDocumentListIcon,
  EllipsisVerticalIcon,
  PencilIcon,
  PencilSquareIcon,
  TrashIcon,
} from "@heroicons/react/24/solid";
import { Button } from "~/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "~/components/ui/dialog";
import { Label } from "~/components/ui/label";
import { Input } from "~/components/ui/input";
import { Link } from "react-router-dom";
import React from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";

export default () => {
  const [data, setData] = useState<Tunnel[]>([]);
  const { toast } = useToast();
  const [tunnelName, setTunnelName] = useState("");
  const onGetTunnels = async () => {
    const { code, data, msg } = await request<Tunnel[]>("/api/tunnels");
    if (code === 0) {
      setData(data);
      return;
    }
    toast({ title: "Failed !", description: msg });
  };

  const onCopyToken = async (id: string) => {
    const resp = await request<string>(`/api/token/${id}`);
    if (resp && resp.code === 0) {
      await navigator.clipboard.writeText(resp.data);
      toast({
        title: "Success !",
        description: "The install token is already copy to clipboard.",
      });
    } else {
      toast({
        title: "Failed !",
        description: resp.msg,
      });
    }
  };

  const onDeleted = async (id: string) => {
    const resp = await request(`/api/tunnels/${id}`, {
      method: "DELETE",
    });
    if (resp && resp.code === 0) {
      onGetTunnels();
      toast({
        title: "Success !",
        description: "tunnel deleted success.",
      });
    } else {
      toast({
        title: "Failed !",
        description: resp.msg,
      });
    }
  };

  const onRefreshToken = async (id: string) => {
    const resp = await request(`/api/tunnels/${id}/refreshtoken`, {
      method: "POST",
    });
    if (resp && resp.code === 0) {
      toast({
        title: "Success !",
        description: "tunnel token refresh success.",
      });
    } else {
      toast({
        title: "Failed !",
        description: resp.msg,
      });
    }
  };

  const handelInputName = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setTunnelName(value);
  };

  const newTunnel = async () => {
    const tunnel: Tunnel = { name: tunnelName };
    const resp = await request("/api/tunnels", {
      method: "POST",
      body: JSON.stringify(tunnel),
    });
    if (resp && resp.code === 0) {
      onGetTunnels();
      toast({
        title: "Success !",
        description: "tunnel add success.",
      });
    } else {
      toast({
        title: "Failed !",
        description: resp.msg,
      });
    }
  };

  useEffect(() => {
    onGetTunnels();
  }, []);

  return (
    <Card className="border-none">
      <CardHeader>
        <CardTitle>Tunnels</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="whitespace-nowrap">Tunnel Name</TableHead>
              <TableHead>Tunnel ID</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Uptime</TableHead>
              <TableHead></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.map((item) => (
              <TableRow key={item.id}>
                <TableCell className="underline">
                  <Link to={`/tunnels/${item.id}`}>{item.name}</Link>
                </TableCell>
                <TableCell className="underline  max-w-60 whitespace-nowrap text-ellipsis overflow-hidden">
                  <Link to={`/tunnels/${item.id}`}>{item.id}</Link>
                </TableCell>
                <TableCell>{item.status}</TableCell>
                <TableCell className="whitespace-nowrap">
                  {item.uptime
                    ? new Date(item.uptime * 1000).toLocaleString()
                    : "--"}
                </TableCell>
                <TableCell className="flex">
                  <Button
                    variant="ghost"
                    className="rounded-full"
                    size="icon"
                    onClick={() => onCopyToken(item.id as string)}
                  >
                    <ClipboardDocumentListIcon className="size-5" />
                  </Button>

                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button
                        variant="ghost"
                        className="rounded-full"
                        size="icon"
                      >
                        <EllipsisVerticalIcon className="size-5" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent>
                      <DropdownMenuItem
                        onClick={() => onRefreshToken(item.id as string)}
                      >
                        <PencilSquareIcon className="h-4 w-4 mr-2" />
                        <span>Edit</span>
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => onRefreshToken(item.id as string)}
                      >
                        <ArrowPathIcon className="h-4 w-4 mr-2" />
                        <span>Refresh Token</span>
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        className="text-destructive"
                        onClick={() => onDeleted(item.id as string)}
                      >
                        <TrashIcon className="h-4 w-4 mr-2" />
                        <span>Deleted</span>
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
      <div className="flex justify-center items-center py-4">
        <Dialog>
          <DialogTrigger asChild>
            <Button className="w-64">Add Tunnel</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>New Tunnel</DialogTitle>
              <DialogDescription>
                Input the name and submit,and then you will get a install token.
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="name" className="text-right">
                  Tunnel Name
                </Label>
                <Input
                  id="name"
                  className="col-span-3"
                  value={tunnelName}
                  onChange={handelInputName}
                />
              </div>
            </div>
            <DialogFooter>
              <DialogClose>
                <Button type="submit" onClick={newTunnel}>
                  Submit
                </Button>
              </DialogClose>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </Card>
  );
};
