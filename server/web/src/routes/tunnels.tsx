import { Card, CardContent } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { request, Tunnel } from "@/lib/request";
import { useEffect, useState } from "react";

export default () => {
  const [data, setData] = useState<Tunnel[]>([]);
  const getTunnels = async () => {
    const data = await request<Tunnel[]>("/api/tunnels");
    setData(data);
  };

  useEffect(() => {
    getTunnels();
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
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.map((item) => (
              <TableRow key={item.id}>
                <TableCell className="underline">{item.name}</TableCell>
                <TableCell>{item.id}</TableCell>
                <TableCell>{item.status}</TableCell>
                <TableCell className="text-right">
                  {new Date(item.uptime * 1000).toLocaleString()}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
};
