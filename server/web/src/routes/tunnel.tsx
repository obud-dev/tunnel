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

export default () => {
  const { id } = useParams();
  const [data, setData] = useState<Route[]>([]);
  const onGetRoutes = async () => {
    const resp = await request<Route[]>(`/api/routes/${id}`);
    if (resp && resp.code === 0) {
      setData(resp.data);
    }
  };

  useEffect(() => {
    onGetRoutes();
  }, []);
  return (
    <Card className="border-none">
      <CardHeader>
        <CardTitle>Routes</CardTitle>
      </CardHeader>
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
