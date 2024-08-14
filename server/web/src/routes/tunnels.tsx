import { Card, CardContent } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useToast } from "@/components/ui/use-toast";
import { request, Tunnel } from "@/lib/request";
import { useEffect, useState } from "react";
import { ClipboardDocumentListIcon } from "@heroicons/react/24/solid";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Link } from "react-router-dom";
import { generateId, generateToken } from "@/lib/utils";
import React from "react";

export default () => {
  const [data, setData] = useState<Tunnel[]>([]);
  const { toast } = useToast();
  const [tunnelName, setTunnelName] = useState("");
  const onGetTunnels = async () => {
    const resp = await request("/api/tunnels");
    if(resp && resp.code===0) {
      setData(resp.data);
    }
  };

  const copyToken = async (id: string) => {
    const resp = await request(`/api/token/${id}`);
    if(resp && resp.code===0) {
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

  const handelInputName = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setTunnelName(value);
  };

  const newTunnel = async () => {
    const tunnel:Tunnel = {
      id: generateId(),
      name: tunnelName,
      token: generateToken(16),
      uptime: Date.now(),
      status: "offline",
    }
    const resp = await request("/api/tunnels",{
      method: "POST",
      body: JSON.stringify(tunnel),
    })
    if(resp && resp.code === 0) {
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
    <Card>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Tunnel ID</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Uptime</TableHead>
              <TableHead>Install Token</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.map((item) => (
              <TableRow key={item.id}>
                <TableCell className="underline">
                  <Link to={`/tunnels/${item.id}`}>{item.name}</Link>
                </TableCell>
                <TableCell className="underline">
                  <Link to={`/tunnels/${item.id}`}>{item.id}</Link>
                </TableCell>
                <TableCell>{item.status}</TableCell>
                <TableCell className="text-right">
                  {new Date(item.uptime * 1000).toLocaleString()}
                </TableCell>
                <TableCell>
                  <ClipboardDocumentListIcon
                    className="size-5 cursor-pointer"
                    onClick={() => copyToken(item.id)}
                  />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
      <Dialog>
        <DialogTrigger>
          <Button>Add Tunnel</Button>
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
    </Card>
  );
};
