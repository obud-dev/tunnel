import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table";
import { request, Route } from "~/lib/request";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "~/components/ui/sheet";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "~/components/ui/use-toast";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select"
import { EllipsisVerticalIcon, PencilSquareIcon, TrashIcon } from "@heroicons/react/24/solid";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "~/components/ui/dropdown-menu";

const FormSchema = z.object({
  id: z.string().optional().readonly(),
  protocol: z.string().min(2, {
    message: "Protocol is too short",
  }),
  hostname: z.string(),
  prefix: z.string(),
  target: z.string(),
  tunnel_id: z.string(),
});

export default () => {
  const { id } = useParams();
  const [data, setData] = useState<Route[]>([]);
  const [open, setOpen] = useState(false);
  const onGetRoutes = async () => {
    const resp = await request<Route[]>(`/api/routes/${id}`);
    if (resp && resp.code === 0) {
      setData(resp.data);
    }
  };

  const defaultValues = {
    id: "",
    tunnel_id: id,
    protocol: "",
    hostname: "",
    prefix: "",
    target: "",
  }
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues,
  });

  const onSubmit = async (data: z.infer<typeof FormSchema>) => {
    setOpen(false)

    const method = data.id ? "PUT" : "POST"
    const { code,msg } = await request("/api/routes",{
      method,
      body: JSON.stringify(data),
    });
    if(code === 0) {
      onGetRoutes();
      toast({
        title: "Success !",
        description: "Route update success.",
      });
      return;
    }
    toast({
      title: "Failed !",
      description: msg,
    });
  }

  const onDelete = async (id: string) => {
    const { code,msg } = await request(`/api/routes/${id}`,{
      method: "DELETE"
    });
    if(code === 0) {
      onGetRoutes();
      toast({
        title: "Success !",
        description: "Route delete success.",
      });
      return;
    }
    toast({
      title: "Failed !",
      description: msg,
    });
  }

  useEffect(() => {
    onGetRoutes();
    if(!open) {
      form.reset(defaultValues)
    }
  }, [open,form]);
  return (
    <Card className="border-none">
      <CardHeader>
        <CardTitle>Routes</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Protocol</TableHead>
              <TableHead>Hostname</TableHead>
              <TableHead>Prefix</TableHead>
              <TableHead className="text-right">Target</TableHead>
              <TableHead></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.map((item) => (
              <TableRow key={item.id}>
                <TableCell>{item.protocol}</TableCell>
                <TableCell className="underline">{item.hostname}</TableCell>
                <TableCell>{item.prefix}</TableCell>
                <TableCell className="text-right">{item.target}</TableCell>
                <TableCell>
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
                        setOpen(true);
                        form.reset(item);
                      }}>
                        <PencilSquareIcon className="h-4 w-4 mr-2" />
                        <span>Edit</span>
                      </DropdownMenuItem>
                      <DropdownMenuItem 
                        className="text-destructive"
                        onClick={() => onDelete(item.id)}
                      >
                        <TrashIcon className="h-4 w-4 mr-2" />
                        <span>Delete</span>
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
        <Sheet open={open} onOpenChange={setOpen}>
          <SheetTrigger asChild>
            <Button className="w-64">Add Route</Button>
          </SheetTrigger>
          <SheetContent aria-describedby={undefined}>
            <SheetHeader>
              <SheetTitle>Add Route</SheetTitle>
            </SheetHeader>
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-6 w-full py-4"
              >
                <FormField
                  control={form.control}
                  name="protocol"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Protocol</FormLabel>
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select a protocol" />
                        </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="http">Http</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="hostname"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Hostname</FormLabel>
                      <FormControl>
                        <Input placeholder="www.example.com" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="prefix"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Prefix</FormLabel>
                      <FormControl>
                        <Input placeholder="/example" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="target"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Target</FormLabel>
                      <FormControl>
                        <Input placeholder="127.0.0.1:8888" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <Button type="submit" className="w-full">
                  Submit
                </Button>
              </form>
            </Form>
          </SheetContent>
        </Sheet>
      </div>
    </Card>
  );
};
