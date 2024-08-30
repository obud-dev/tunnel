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
  PencilSquareIcon,
  TrashIcon,
} from "@heroicons/react/24/solid";
import { Button } from "~/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "~/components/ui/dialog";
import { Input } from "~/components/ui/input";
import { Link } from "react-router-dom";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "~/components/ui/form";
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from "~/components/ui/sheet";
import { uptime } from "process";

export default () => {
  const [data, setData] = useState<Tunnel[]>([]);
  const [open, setOpen] = useState(false);
  const [openEdit, setOpenEdit] = useState(false);
  const { toast } = useToast();
  const FormSchema = z.object({
    name: z.string().min(4,{
      message: "Tunnel name is too short",
    }),
    id: z.string().optional().readonly(),
    token: z.string().optional().readonly(),
    status: z.string().optional().readonly(),
  })

  const defaultValues = {
    id: "",
    name: "",
    token: "",
    status: "",
  }

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
  })

  const onGetTunnels = async () => {
    const { code, data, msg } = await request<Tunnel[]>("/api/tunnels");
    if (code === 0) {
      setData(data);
      return;
    }
    toast({ title: "Failed !", description: msg });
  };

  const onCopyToken = async (id: string) => {
    const { code,data,msg } = await request<string>(`/api/token/${id}`);
    if (code === 0) {
      await navigator.clipboard.writeText(data);
      toast({
        title: "Success !",
        description: "The install token is already copy to clipboard.",
      });
      return;
    }
    toast({
      title: "Failed !",
      description: msg,
    });
  };

  const onDeleted = async (id: string) => {
    const { code,msg } = await request(`/api/tunnels/${id}`, {
      method: "DELETE",
    });
    if (code === 0) {
      onGetTunnels();
      toast({
        title: "Success !",
        description: "tunnel deleted success.",
      });
      return;
    }
    toast({
      title: "Failed !",
      description: msg,
    });
  };

  const onRefreshToken = async (id: string) => {
    const { code,msg } = await request(`/api/tunnels/${id}/refreshtoken`, {
      method: "POST",
    });
    if (code === 0) {
      toast({
        title: "Success !",
        description: "tunnel token refresh success.",
      });
      return;
    }
    toast({
      title: "Failed !",
      description: msg
    });
  };

  const newTunnel = async (data: z.infer<typeof FormSchema>) => {
    setOpen(false)
    const { code,msg } = await request("/api/tunnels", {
      method: "POST",
      body: JSON.stringify(data),
    });
    if (code === 0) {
      onGetTunnels();
      toast({
        title: "Success !",
        description: "tunnel add success.",
      });
      return;
    }
    toast({
      title: "Failed !",
      description: msg,
    });
  };

  const onEditSubmit = async (data: z.infer<typeof FormSchema>) => {
    setOpenEdit(false)
    const { code,msg } = await request(`/api/tunnels/${data.id}`, {
      method: "PUT",
      body: JSON.stringify(data)
    })
    if (code === 0) {
      onGetTunnels();
      toast({
        title: "Success !",
        description: "tunnel edit success.",
      });
      return;
    }
    toast({
      title: "Failed !",
      description: msg,
    });
  }

  useEffect(() => {
    onGetTunnels();
    if(!openEdit) {
      form.reset(defaultValues)
    }
  }, [openEdit,form]);

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
                      <DropdownMenuItem onClick={() => {
                        setOpenEdit(true);
                        form.reset(item);
                      }}>
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
        <Dialog open={open} onOpenChange={setOpen}>
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
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(newTunnel)}
                className="space-y-6 w-full"
              >
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Tunnel Name</FormLabel>
                      <FormControl>
                        <Input {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <Button type="submit" className="w-full">Submit</Button>
              </form>
            </Form>
          </DialogContent>
        </Dialog>
        <Sheet open={openEdit} onOpenChange={setOpenEdit}>
          <SheetContent>
            <SheetHeader>
              <SheetTitle>Edit Tunnel</SheetTitle>
              <SheetDescription>
                {form.getValues().id}
              </SheetDescription>
            </SheetHeader>
            <Form {...form}>
              <form
                  onSubmit={form.handleSubmit(onEditSubmit)}
                  className="space-y-6 w-full py-4"
                >
                <FormField 
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Tunnel Name</FormLabel>
                      <FormControl>
                        <Input {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <Button type="submit" className="w-full">
                  Update
                </Button>
              </form>
            </Form>
          </SheetContent>
        </Sheet>
      </div>
    </Card>
  );
};
