#ifndef GOLANG_INTERFACE_H
#define GOLANG_INTERFACE_H

#ifdef __cplusplus
extern "C" {
#endif
  //! This definitions is used to interface the golang counterpart.
  typedef struct PublicAttributes
  {
      unsigned long  id;
      unsigned       entityType;
      double         size;
      long           posX;
      long           posY;
      double         thrustR;
      double         thrustTheta;
  }PublicAttributes;

  typedef struct WorldConfig
  {
    unsigned long foodCount;
    unsigned foodProductionRate;
    unsigned long creatureCount;
  } WorldConfig;

  //! Generates the world based on the input config.
  void Generate(WorldConfig config);
  //! Runs the engine once.
  void Execute();
  //! Returns the raw pointer of entity public attributes vector.
  PublicAttributes GetPublicAttribute(int index);
  //! Return the number of entities in the world space.
  int GetEntityCount();
  //! Sets thrust values of the selected entity.
  void SetThrust(int id, double r, double theta);

#ifdef __cplusplus
}
#endif

#endif // GOLANG_INTERFACE_H