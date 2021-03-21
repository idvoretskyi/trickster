/*
 * Copyright 2018 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dataset

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Hash) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 uint64
		zb0001, err = dc.ReadUint64()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = Hash(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Hash) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteUint64(uint64(z))
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Hash) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendUint64(o, uint64(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Hash) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 uint64
		zb0001, bts, err = msgp.ReadUint64Bytes(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = Hash(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Hash) Msgsize() (s int) {
	s = msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Hashes) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(Hashes, zb0002)
	}
	for zb0001 := range *z {
		{
			var zb0003 uint64
			zb0003, err = dc.ReadUint64()
			if err != nil {
				err = msgp.WrapError(err, zb0001)
				return
			}
			(*z)[zb0001] = Hash(zb0003)
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Hashes) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0004 := range z {
		err = en.WriteUint64(uint64(z[zb0004]))
		if err != nil {
			err = msgp.WrapError(err, zb0004)
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Hashes) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zb0004 := range z {
		o = msgp.AppendUint64(o, uint64(z[zb0004]))
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Hashes) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(Hashes, zb0002)
	}
	for zb0001 := range *z {
		{
			var zb0003 uint64
			zb0003, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, zb0001)
				return
			}
			(*z)[zb0001] = Hash(zb0003)
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Hashes) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize + (len(z) * (msgp.Uint64Size))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Series) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "header":
			err = z.Header.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Header")
				return
			}
		case "points":
			err = z.Points.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Points")
				return
			}
		case "ps":
			z.PointSize, err = dc.ReadInt64()
			if err != nil {
				err = msgp.WrapError(err, "PointSize")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Series) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "header"
	err = en.Append(0x83, 0xa6, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72)
	if err != nil {
		return
	}
	err = z.Header.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Header")
		return
	}
	// write "points"
	err = en.Append(0xa6, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73)
	if err != nil {
		return
	}
	err = z.Points.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Points")
		return
	}
	// write "ps"
	err = en.Append(0xa2, 0x70, 0x73)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.PointSize)
	if err != nil {
		err = msgp.WrapError(err, "PointSize")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Series) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "header"
	o = append(o, 0x83, 0xa6, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72)
	o, err = z.Header.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Header")
		return
	}
	// string "points"
	o = append(o, 0xa6, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73)
	o, err = z.Points.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Points")
		return
	}
	// string "ps"
	o = append(o, 0xa2, 0x70, 0x73)
	o = msgp.AppendInt64(o, z.PointSize)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Series) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "header":
			bts, err = z.Header.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Header")
				return
			}
		case "points":
			bts, err = z.Points.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Points")
				return
			}
		case "ps":
			z.PointSize, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "PointSize")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Series) Msgsize() (s int) {
	s = 1 + 7 + z.Header.Msgsize() + 7 + z.Points.Msgsize() + 3 + msgp.Int64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *SeriesLookupKey) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "StatementID":
			z.StatementID, err = dc.ReadInt()
			if err != nil {
				err = msgp.WrapError(err, "StatementID")
				return
			}
		case "Hash":
			{
				var zb0002 uint64
				zb0002, err = dc.ReadUint64()
				if err != nil {
					err = msgp.WrapError(err, "Hash")
					return
				}
				z.Hash = Hash(zb0002)
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z SeriesLookupKey) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "StatementID"
	err = en.Append(0x82, 0xab, 0x53, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44)
	if err != nil {
		return
	}
	err = en.WriteInt(z.StatementID)
	if err != nil {
		err = msgp.WrapError(err, "StatementID")
		return
	}
	// write "Hash"
	err = en.Append(0xa4, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = en.WriteUint64(uint64(z.Hash))
	if err != nil {
		err = msgp.WrapError(err, "Hash")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z SeriesLookupKey) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "StatementID"
	o = append(o, 0x82, 0xab, 0x53, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44)
	o = msgp.AppendInt(o, z.StatementID)
	// string "Hash"
	o = append(o, 0xa4, 0x48, 0x61, 0x73, 0x68)
	o = msgp.AppendUint64(o, uint64(z.Hash))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *SeriesLookupKey) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "StatementID":
			z.StatementID, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "StatementID")
				return
			}
		case "Hash":
			{
				var zb0002 uint64
				zb0002, bts, err = msgp.ReadUint64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Hash")
					return
				}
				z.Hash = Hash(zb0002)
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z SeriesLookupKey) Msgsize() (s int) {
	s = 1 + 12 + msgp.IntSize + 5 + msgp.Uint64Size
	return
}