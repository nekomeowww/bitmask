package bitmask

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBitMask(t *testing.T) {
	assert := assert.New(t)

	var b BitMask
	b.Set(1)
	assert.Equal(1, int(b))
	b.Set(2)
	assert.Equal(1+2, int(b))
	b.Set(3)
	assert.Equal(1+2+4, int(b))
	assert.True(b.IsSet(1))
	assert.True(b.IsSet(2))
	assert.True(b.IsSet(3))
	assert.False(b.IsSet(4))

	b.Unset(1)
	assert.Equal(2+4, int(b))

	b.Unset(2)
	assert.Equal(4, int(b))

	b.Unset(3)
	assert.Equal(0, int(b))

	assert.False(b.IsSet(16))
	b.Set(16)
	assert.Equal(uint32(1<<15), uint32(b))
	assert.True(b.IsSet(16))

	b.Unset(16)
	assert.Equal(0, int(b))
	assert.False(b.IsSet(16))

	assert.False(b.IsSet(32))
	b.Set(32)
	assert.Equal(uint32(1<<31), uint32(b))
	assert.True(b.IsSet(32))

	b.Unset(32)
	assert.Equal(0, int(b))
	assert.False(b.IsSet(32))
}

func TestMarshal(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		t.Run("ErrNotAPointerOrStruct", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			mData := map[string]interface{}{}
			bitmask, err := Marshal(mData)
			require.Error(err)
			require.NotNil(bitmask)

			assert.ErrorIs(err, ErrNotStructPointerOrStruct)

			bitmask, err = Marshal(&mData)
			require.Error(err)
			require.NotNil(bitmask)

			assert.ErrorIs(err, ErrNotStructPointerOrStruct)
		})

		t.Run("ErrInvalidTagValue", func(t *testing.T) {
			type tS1 struct {
				A bool `bitmask:"one"`
				B bool `bitmask:"two"`
			}

			assert := assert.New(t)
			require := require.New(t)

			bitmask, err := Marshal(tS1{})
			require.Error(err)
			require.NotNil(bitmask)

			assert.ErrorIs(err, ErrInvalidTagValue)
		})

		t.Run("ErrBitOutOfRange", func(t *testing.T) {
			t.Run("Beyond1", func(t *testing.T) {
				type tS1 struct {
					A bool `bitmask:"0"`
					B bool `bitmask:"1"`
				}

				assert := assert.New(t)
				require := require.New(t)

				bitmask, err := Marshal(tS1{})
				require.Error(err)
				require.NotNil(bitmask)

				assert.ErrorIs(err, ErrBitOutOfRange)
			})

			t.Run("Beyond32", func(t *testing.T) {
				type tS1 struct {
					A bool `bitmask:"1"`
					B bool `bitmask:"65"`
				}

				assert := assert.New(t)
				require := require.New(t)

				bitmask, err := Marshal(tS1{})
				require.Error(err)
				require.NotNil(bitmask)

				assert.ErrorIs(err, ErrBitOutOfRange)
			})
		})
	})

	t.Run("Nil", func(t *testing.T) {
		t.Run("Interface", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			bitmask, err := Marshal(nil)
			require.NoError(err)
			require.NotNil(bitmask)

			assert.Equal(0, int(bitmask))
		})

		t.Run("Pointer", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			type tS1 struct {
				A bool `bitmask:"1"`
				B bool `bitmask:"2"`
			}

			var val *tS1
			bitmask, err := Marshal(val)
			require.NoError(err)
			require.NotNil(bitmask)

			assert.Equal(0, int(bitmask))
		})
	})

	t.Run("EmptyTag", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		type tS1 struct {
			A bool `bitmask:""`
			B bool `bitmask:"2"`
		}

		bitmask, err := Marshal(tS1{A: true, B: true})
		require.NoError(err)
		require.NotNil(bitmask)

		assert.Equal(2, int(bitmask))
		assert.False(bitmask.IsSet(1))
		assert.True(bitmask.IsSet(2))
	})

	t.Run("Bool", func(t *testing.T) {
		type tS1 struct {
			A bool `bitmask:"1"`
			B bool `bitmask:"2"`
		}

		t.Run("Struct", func(t *testing.T) {
			t.Run("AllSet", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				v := tS1{A: true, B: true}
				bitmask, err := Marshal(v)
				require.NoError(err)
				require.NotNil(bitmask)

				assert.True(bitmask.IsSet(1))
				assert.True(bitmask.IsSet(2))
				assert.Equal(3, int(bitmask))
			})

			t.Run("Partial", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				v := tS1{A: true}
				bitmask, err := Marshal(v)
				require.NoError(err)
				require.NotNil(bitmask)

				assert.True(bitmask.IsSet(1))
				assert.False(bitmask.IsSet(2))
				assert.Equal(1, int(bitmask))
			})
		})

		t.Run("StructPointer", func(t *testing.T) {
			t.Run("AllSet", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				v := tS1{A: true, B: true}
				bitmask, err := Marshal(&v)
				require.NoError(err)
				require.NotNil(bitmask)

				assert.True(bitmask.IsSet(1))
				assert.True(bitmask.IsSet(2))
				assert.Equal(3, int(bitmask))
			})

			t.Run("Partial", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				v := tS1{A: true}
				bitmask, err := Marshal(&v)
				require.NoError(err)
				require.NotNil(bitmask)

				assert.True(bitmask.IsSet(1))
				assert.False(bitmask.IsSet(2))
				assert.Equal(1, int(bitmask))
			})
		})
	})

	t.Run("Int", func(t *testing.T) {
		type tS1 struct {
			A int `bitmask:"1"`
			B int `bitmask:"2"`
		}

		t.Run("Struct", func(t *testing.T) {
			t.Run("AllSet", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				v := tS1{A: 1, B: 1}
				bitmask, err := Marshal(v)
				require.NoError(err)
				require.NotNil(bitmask)

				assert.True(bitmask.IsSet(1))
				assert.True(bitmask.IsSet(2))
				assert.Equal(3, int(bitmask))
			})

			t.Run("Partial", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				v := tS1{A: 1}
				bitmask, err := Marshal(v)
				require.NoError(err)
				require.NotNil(bitmask)

				assert.True(bitmask.IsSet(1))
				assert.False(bitmask.IsSet(2))
				assert.Equal(1, int(bitmask))
			})
		})

		t.Run("StructPointer", func(t *testing.T) {
			t.Run("AllSet", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				v := tS1{A: 1, B: 1}
				bitmask, err := Marshal(&v)
				require.NoError(err)
				require.NotNil(bitmask)

				assert.True(bitmask.IsSet(1))
				assert.True(bitmask.IsSet(2))
				assert.Equal(3, int(bitmask))
			})

			t.Run("Partial", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				v := tS1{A: 1}
				bitmask, err := Marshal(&v)
				require.NoError(err)
				require.NotNil(bitmask)

				assert.True(bitmask.IsSet(1))
				assert.False(bitmask.IsSet(2))
				assert.Equal(1, int(bitmask))
			})
		})
	})

	t.Run("Uint", func(t *testing.T) {
		type tS1 struct {
			A uint `bitmask:"1"`
			B uint `bitmask:"2"`
		}

		t.Run("Struct", func(t *testing.T) {
			t.Run("AllSet", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				v := tS1{A: 1, B: 1}
				bitmask, err := Marshal(v)
				require.NoError(err)
				require.NotNil(bitmask)

				assert.True(bitmask.IsSet(1))
				assert.True(bitmask.IsSet(2))
				assert.Equal(3, int(bitmask))
			})

			t.Run("Partial", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				v := tS1{A: 1}
				bitmask, err := Marshal(v)
				require.NoError(err)
				require.NotNil(bitmask)

				assert.True(bitmask.IsSet(1))
				assert.False(bitmask.IsSet(2))
				assert.Equal(1, int(bitmask))
			})
		})

		t.Run("StructPointer", func(t *testing.T) {
			t.Run("AllSet", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				v := tS1{A: 1, B: 1}
				bitmask, err := Marshal(&v)
				require.NoError(err)
				require.NotNil(bitmask)

				assert.True(bitmask.IsSet(1))
				assert.True(bitmask.IsSet(2))
				assert.Equal(3, int(bitmask))
			})

			t.Run("Partial", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				v := tS1{A: 1}
				bitmask, err := Marshal(&v)
				require.NoError(err)
				require.NotNil(bitmask)

				assert.True(bitmask.IsSet(1))
				assert.False(bitmask.IsSet(2))
				assert.Equal(1, int(bitmask))
			})
		})
	})
}

func TestUnmarshal(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		t.Run("ErrNotAStructPointer", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			bitmask := New(3)
			mData := map[string]interface{}{}

			err := Unmarshal(bitmask, &mData)
			require.Error(err)
			require.NotNil(bitmask)

			assert.ErrorIs(err, ErrNotAStructPointer)
		})

		t.Run("ErrInvalidTagValue", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			type tS1 struct {
				A bool `bitmask:"one"`
				B bool `bitmask:"two"`
			}

			bitmask := New(3)
			err := Unmarshal(bitmask, &tS1{})
			require.Error(err)

			assert.ErrorIs(err, ErrInvalidTagValue)
		})

		t.Run("ErrBitOutOfRange", func(t *testing.T) {
			t.Run("Beyond1", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				type tS1 struct {
					A bool `bitmask:"0"`
					B bool `bitmask:"1"`
				}

				bitmask := New(3)
				err := Unmarshal(bitmask, &tS1{})
				require.Error(err)

				assert.ErrorIs(err, ErrBitOutOfRange)
			})

			t.Run("Beyond32", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				type tS1 struct {
					A bool `bitmask:"1"`
					B bool `bitmask:"65"`
				}

				bitmask := New(3)
				err := Unmarshal(bitmask, &tS1{})
				require.Error(err)

				assert.ErrorIs(err, ErrBitOutOfRange)
			})
		})
	})

	t.Run("Nil", func(t *testing.T) {
		t.Run("Interface", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			bitmask := New(3)
			err := Unmarshal(bitmask, nil)
			require.Error(err)

			assert.ErrorIs(err, ErrNotAStructPointer)
		})

		t.Run("Pointer", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			bitmask := New(3)
			err := Unmarshal(bitmask, (*int)(nil))
			require.Error(err)

			assert.ErrorIs(err, ErrNotAStructPointer)
		})
	})

	t.Run("EmptyTag", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		type tS1 struct {
			A bool `bitmask:""`
			B bool `bitmask:"2"`
		}

		var v tS1
		bitmask := New(3)
		assert.True(bitmask.IsSet(1))
		assert.True(bitmask.IsSet(2))

		err := Unmarshal(bitmask, &v)
		require.NoError(err)
		require.NotNil(bitmask)

		assert.False(v.A)
		assert.True(v.B)
	})

	t.Run("Bool", func(t *testing.T) {
		type tS1 struct {
			A bool `bitmask:"1"`
			B bool `bitmask:"2"`
		}

		t.Run("Struct", func(t *testing.T) {
			t.Run("AllSet", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				var v tS1
				bitmask := New(3)
				assert.True(bitmask.IsSet(1))
				assert.True(bitmask.IsSet(2))

				err := Unmarshal(bitmask, &v)
				require.NoError(err)

				assert.True(v.A)
				assert.True(v.B)
			})

			t.Run("Partial", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				var v tS1
				bitmask := New(1)
				assert.True(bitmask.IsSet(1))

				err := Unmarshal(bitmask, &v)
				require.NoError(err)

				assert.True(v.A)
				assert.False(v.B)
			})
		})
	})

	t.Run("Int", func(t *testing.T) {
		type tS1 struct {
			A int `bitmask:"1"`
			B int `bitmask:"2"`
		}

		t.Run("Struct", func(t *testing.T) {
			t.Run("AllSet", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				var v tS1
				bitmask := New(3)
				assert.True(bitmask.IsSet(1))
				assert.True(bitmask.IsSet(2))

				err := Unmarshal(bitmask, &v)
				require.NoError(err)

				assert.Equal(1, v.A)
				assert.Equal(1, v.B)
			})

			t.Run("Partial", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				var v tS1
				bitmask := New(1)
				assert.True(bitmask.IsSet(1))

				err := Unmarshal(bitmask, &v)
				require.NoError(err)

				assert.Equal(1, v.A)
				assert.Equal(0, v.B)
			})
		})
	})

	t.Run("Uint", func(t *testing.T) {
		type tS1 struct {
			A uint `bitmask:"1"`
			B uint `bitmask:"2"`
		}

		t.Run("Struct", func(t *testing.T) {
			t.Run("AllSet", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				var v tS1
				bitmask := New(3)
				assert.True(bitmask.IsSet(1))
				assert.True(bitmask.IsSet(2))

				err := Unmarshal(bitmask, &v)
				require.NoError(err)

				assert.Equal(uint(1), v.A)
				assert.Equal(uint(1), v.B)
			})

			t.Run("Partial", func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				var v tS1
				bitmask := New(1)
				assert.True(bitmask.IsSet(1))

				err := Unmarshal(bitmask, &v)
				require.NoError(err)

				assert.Equal(uint(1), v.A)
				assert.Equal(uint(0), v.B)
			})
		})
	})
}
