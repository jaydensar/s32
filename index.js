import "dotenv/config";

import express from "express";
import Enmap from "enmap";

const data = new Enmap({ name: "s32" });

const app = express();

app.use(express.json());

const keys = {
  GET: process.env.GET_KEY,
  PUT: process.env.PUT_KEY,
};

app.get("/game", (req, res) => {
  if (req.query.Key !== keys.GET || !req.query.UserID) {
    return res.sendStatus(401);
  }

  console.log("got valid get request", req.query, req.headers);

  data.ensure(req.query.UserID, {
    Challenges: {},
    Inventory: {},
    Points: 0,
  });

  res.status(200).json(data.get(req.query.UserID));
});

app.put("/game", (req, res) => {
  if (req.body.Key !== keys.PUT || !req.body.UserID) {
    return res.sendStatus(401);
  }

  console.log("got valid put request", req.body, req.headers);

  data.set(String(req.body.UserID), {
    Challenges: req.body.Challenges,
    Inventory: req.body.Inventory,
    Points: req.body.Points,
  });

  res.sendStatus(204);
});

app.listen(process.env.PORT || 3000, () => {
  console.log(`Listening on port ${process.env.port || 3000}`);
});
