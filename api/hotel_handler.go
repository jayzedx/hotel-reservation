package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
)

// struct
type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

// function
func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

// implement
type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	// fmt.Println(qparams.Rooms, qparams.Rating)

	query := bson.M{}
	if c.Query("rooms") != "" {
		query["rooms"] = bson.M{"$exists": qparams.Rooms}
	}
	if c.Query("rating") != "" {
		query["rating"] = qparams.Rating
	}
	hotels, err := h.hotelStore.GetHotels(c.Context(), query)

	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error {

	return nil
}

func (h *HotelHandler) HandlePutHotel(c *fiber.Ctx) error {

	return nil
}

func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {

	return nil
}
