//go:build !FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build !FEATURE_TOGGLE_MARGIN_CLI_ALPHA

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sifnode/clp/v1/pool.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Pool struct {
	ExternalAsset                 *Asset                                  `protobuf:"bytes,1,opt,name=external_asset,json=externalAsset,proto3" json:"external_asset,omitempty"`
	NativeAssetBalance            github_com_cosmos_cosmos_sdk_types.Uint `protobuf:"bytes,2,opt,name=native_asset_balance,json=nativeAssetBalance,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"native_asset_balance" yaml:"native_asset_balance"`
	ExternalAssetBalance          github_com_cosmos_cosmos_sdk_types.Uint `protobuf:"bytes,3,opt,name=external_asset_balance,json=externalAssetBalance,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"external_asset_balance" yaml:"external_asset_balance"`
	PoolUnits                     github_com_cosmos_cosmos_sdk_types.Uint `protobuf:"bytes,4,opt,name=pool_units,json=poolUnits,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"pool_units" yaml:"pool_units"`
	SwapPriceNative               *github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=swap_price_native,json=swapPriceNative,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"swap_price_native,omitempty" yaml:"swap_price_native "`
	SwapPriceExternal             *github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,6,opt,name=swap_price_external,json=swapPriceExternal,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"swap_price_external,omitempty" yaml:"swap_price_external "`
	RewardPeriodNativeDistributed github_com_cosmos_cosmos_sdk_types.Uint `protobuf:"bytes,7,opt,name=reward_period_native_distributed,json=rewardPeriodNativeDistributed,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"reward_period_native_distributed" yaml:"reward_period_native_distributed"`
}

func (m *Pool) Reset()         { *m = Pool{} }
func (m *Pool) String() string { return proto.CompactTextString(m) }
func (*Pool) ProtoMessage()    {}
func (*Pool) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf43b9338450af59, []int{0}
}
func (m *Pool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Pool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Pool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Pool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pool.Merge(m, src)
}
func (m *Pool) XXX_Size() int {
	return m.Size()
}
func (m *Pool) XXX_DiscardUnknown() {
	xxx_messageInfo_Pool.DiscardUnknown(m)
}

var xxx_messageInfo_Pool proto.InternalMessageInfo

func (m *Pool) GetExternalAsset() *Asset {
	if m != nil {
		return m.ExternalAsset
	}
	return nil
}

func init() {
	proto.RegisterType((*Pool)(nil), "sifnode.clp.v1.Pool")
}

func init() { proto.RegisterFile("sifnode/clp/v1/pool.proto", fileDescriptor_cf43b9338450af59) }

var fileDescriptor_cf43b9338450af59 = []byte{
	// 471 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0xc1, 0x6e, 0xd3, 0x30,
	0x18, 0xc7, 0x6b, 0x18, 0x45, 0x33, 0x62, 0xa8, 0xa6, 0xa0, 0xac, 0xd3, 0x92, 0xc8, 0x07, 0xb6,
	0x0b, 0x89, 0x06, 0x37, 0xc4, 0x65, 0x55, 0x11, 0xe2, 0x32, 0x55, 0x41, 0xbb, 0xec, 0x12, 0xb9,
	0x8e, 0xd7, 0x59, 0xa4, 0xb1, 0x15, 0xbb, 0xdd, 0x7a, 0xe3, 0xc4, 0x99, 0x77, 0xe0, 0x65, 0x26,
	0x71, 0xd9, 0x11, 0xed, 0x10, 0xa1, 0xf6, 0x0d, 0xf6, 0x04, 0x28, 0x76, 0xb2, 0xad, 0x6c, 0x12,
	0xe4, 0x54, 0xfb, 0xff, 0x7d, 0xfd, 0xff, 0xfe, 0xfe, 0xa2, 0x0f, 0x6e, 0x2a, 0x7e, 0x9c, 0x89,
	0x84, 0x85, 0x34, 0x95, 0xe1, 0x6c, 0x2f, 0x94, 0x42, 0xa4, 0x81, 0xcc, 0x85, 0x16, 0x68, 0xa3,
	0x2a, 0x05, 0x34, 0x95, 0xc1, 0x6c, 0xaf, 0xd7, 0x1d, 0x8b, 0xb1, 0x30, 0xa5, 0xb0, 0x3c, 0xd9,
	0xae, 0x5e, 0xef, 0x2f, 0x03, 0x3d, 0x97, 0x4c, 0xd9, 0x1a, 0xfe, 0xd9, 0x86, 0x6b, 0x43, 0x21,
	0x52, 0xf4, 0x1e, 0x6e, 0xb0, 0x33, 0xcd, 0xf2, 0x8c, 0xa4, 0x31, 0x51, 0x8a, 0x69, 0x07, 0xf8,
	0x60, 0xf7, 0xc9, 0x9b, 0x17, 0xc1, 0x2a, 0x23, 0xd8, 0x2f, 0x8b, 0xd1, 0xd3, 0xba, 0xd9, 0x5c,
	0xd1, 0x57, 0x00, 0xbb, 0x19, 0xd1, 0x7c, 0xc6, 0xec, 0x9f, 0xe3, 0x11, 0x49, 0x49, 0x46, 0x99,
	0xf3, 0xc0, 0x07, 0xbb, 0xeb, 0xfd, 0x83, 0xf3, 0xc2, 0x6b, 0x5d, 0x16, 0xde, 0xce, 0x98, 0xeb,
	0x93, 0xe9, 0x28, 0xa0, 0x62, 0x12, 0x52, 0xa1, 0x26, 0x42, 0x55, 0x3f, 0xaf, 0x55, 0xf2, 0xa5,
	0xca, 0x75, 0xc8, 0x33, 0x7d, 0x55, 0x78, 0x5b, 0x73, 0x32, 0x49, 0xdf, 0xe1, 0xfb, 0x4c, 0x71,
	0x84, 0xac, 0x6c, 0xd8, 0x7d, 0x2b, 0xa2, 0x6f, 0x00, 0xbe, 0x5c, 0x7d, 0xc1, 0x75, 0x88, 0x87,
	0x26, 0xc4, 0xb0, 0x79, 0x88, 0x6d, 0x1b, 0xe2, 0x7e, 0x5b, 0x1c, 0x75, 0x57, 0x86, 0x50, 0x07,
	0xa1, 0x10, 0x96, 0x9f, 0x28, 0x9e, 0x66, 0x5c, 0x2b, 0x67, 0xcd, 0xb0, 0x07, 0xcd, 0xd9, 0x1d,
	0xcb, 0xbe, 0xb1, 0xc2, 0xd1, 0x7a, 0x79, 0x39, 0x2c, 0xcf, 0x48, 0xc1, 0x8e, 0x3a, 0x25, 0x32,
	0x96, 0x39, 0xa7, 0x2c, 0xb6, 0xe3, 0x70, 0x1e, 0x19, 0xd6, 0xc7, 0xcb, 0xc2, 0x7b, 0xf5, 0x1f,
	0x9c, 0x01, 0xa3, 0x57, 0x85, 0xb7, 0x69, 0x31, 0x77, 0xcc, 0x7c, 0x1c, 0x3d, 0x2b, 0xc5, 0x61,
	0xa9, 0x1d, 0x18, 0x09, 0xcd, 0xe1, 0xf3, 0x5b, 0x7d, 0xf5, 0xe3, 0x9d, 0xb6, 0xc1, 0x7e, 0x6a,
	0x84, 0xdd, 0xba, 0x83, 0xad, 0xed, 0x7c, 0x1c, 0x75, 0xae, 0xc1, 0x1f, 0x2a, 0x11, 0xfd, 0x00,
	0xd0, 0xcf, 0xd9, 0x29, 0xc9, 0x93, 0x58, 0xb2, 0x9c, 0x8b, 0xa4, 0x8a, 0x19, 0x27, 0x5c, 0xe9,
	0x9c, 0x8f, 0xa6, 0x9a, 0x25, 0xce, 0x63, 0x13, 0xe4, 0xa8, 0xf9, 0xac, 0x77, 0x6c, 0x9a, 0x7f,
	0x01, 0x70, 0xb4, 0x6d, 0x5b, 0x86, 0xa6, 0xc3, 0x4e, 0x65, 0x70, 0x53, 0xef, 0xef, 0x9f, 0x2f,
	0x5c, 0x70, 0xb1, 0x70, 0xc1, 0xef, 0x85, 0x0b, 0xbe, 0x2f, 0xdd, 0xd6, 0xc5, 0xd2, 0x6d, 0xfd,
	0x5a, 0xba, 0xad, 0xa3, 0xdb, 0x61, 0x3e, 0xf3, 0x63, 0x7a, 0x42, 0x78, 0x16, 0xd6, 0x7b, 0x79,
	0x66, 0x36, 0xd3, 0x24, 0x1a, 0xb5, 0xcd, 0x5e, 0xbe, 0xfd, 0x13, 0x00, 0x00, 0xff, 0xff, 0x38,
	0x5f, 0x1b, 0x0f, 0xf6, 0x03, 0x00, 0x00,
}

func (m *Pool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Pool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Pool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.RewardPeriodNativeDistributed.Size()
		i -= size
		if _, err := m.RewardPeriodNativeDistributed.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	if m.SwapPriceExternal != nil {
		{
			size := m.SwapPriceExternal.Size()
			i -= size
			if _, err := m.SwapPriceExternal.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintPool(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.SwapPriceNative != nil {
		{
			size := m.SwapPriceNative.Size()
			i -= size
			if _, err := m.SwapPriceNative.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintPool(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	{
		size := m.PoolUnits.Size()
		i -= size
		if _, err := m.PoolUnits.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.ExternalAssetBalance.Size()
		i -= size
		if _, err := m.ExternalAssetBalance.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.NativeAssetBalance.Size()
		i -= size
		if _, err := m.NativeAssetBalance.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.ExternalAsset != nil {
		{
			size, err := m.ExternalAsset.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPool(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintPool(dAtA []byte, offset int, v uint64) int {
	offset -= sovPool(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Pool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ExternalAsset != nil {
		l = m.ExternalAsset.Size()
		n += 1 + l + sovPool(uint64(l))
	}
	l = m.NativeAssetBalance.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.ExternalAssetBalance.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.PoolUnits.Size()
	n += 1 + l + sovPool(uint64(l))
	if m.SwapPriceNative != nil {
		l = m.SwapPriceNative.Size()
		n += 1 + l + sovPool(uint64(l))
	}
	if m.SwapPriceExternal != nil {
		l = m.SwapPriceExternal.Size()
		n += 1 + l + sovPool(uint64(l))
	}
	l = m.RewardPeriodNativeDistributed.Size()
	n += 1 + l + sovPool(uint64(l))
	return n
}

func sovPool(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPool(x uint64) (n int) {
	return sovPool(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Pool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPool
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Pool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Pool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExternalAsset", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ExternalAsset == nil {
				m.ExternalAsset = &Asset{}
			}
			if err := m.ExternalAsset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NativeAssetBalance", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.NativeAssetBalance.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExternalAssetBalance", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ExternalAssetBalance.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolUnits", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.PoolUnits.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SwapPriceNative", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Dec
			m.SwapPriceNative = &v
			if err := m.SwapPriceNative.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SwapPriceExternal", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Dec
			m.SwapPriceExternal = &v
			if err := m.SwapPriceExternal.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardPeriodNativeDistributed", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RewardPeriodNativeDistributed.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPool(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPool
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipPool(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPool
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPool
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPool
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthPool
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPool
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPool
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPool        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPool          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPool = fmt.Errorf("proto: unexpected end of group")
)
