import { Card, CardContent } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { request, Route } from "@/lib/request";
import { useEffect, useState } from "react";

export default () => {
  const [data, setData] = useState<Route[]>([]);
  const getTunnels = async () => {
    const data = await request<Route[]>("/api/routes");
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
              <TableHead>Hostname</TableHead>
              <TableHead>Prefix</TableHead>
              <TableHead>Protocol</TableHead>
              <TableHead className="text-right">Target</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.map((item) => (
              <TableRow key={item.id}>
                <TableCell className="underline">{item.hostname}</TableCell>
                <TableCell>{item.prefix}</TableCell>
                <TableCell>{item.protocol}</TableCell>
                <TableCell className="text-right">{item.target}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
};
